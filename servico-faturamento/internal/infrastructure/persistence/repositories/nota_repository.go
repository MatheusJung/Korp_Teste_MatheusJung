package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nf-system/servico-faturamento/internal/application/ports"
	"github.com/nf-system/servico-faturamento/internal/domain/entities"
	domainerrors "github.com/nf-system/servico-faturamento/internal/domain/errors"
)

type NotaFiscalRepository struct {
	db *pgxpool.Pool
}

func NewNotaFiscalRepository(db *pgxpool.Pool) *NotaFiscalRepository {
	return &NotaFiscalRepository{db: db}
}

func (r *NotaFiscalRepository) ProximoNumero(ctx context.Context) (int64, error) {
	var numero int64
	err := r.db.QueryRow(ctx, "SELECT nextval('nota_fiscal_numero_seq')").Scan(&numero)
	return numero, err
}

func (r *NotaFiscalRepository) Salvar(ctx context.Context, nota *entities.NotaFiscal) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO notas_fiscais (id, numero, status, criado_em, atualizado_em)
		VALUES ($1, $2, $3, $4, $5)`,
		nota.ID, nota.Numero, nota.Status, nota.CriadoEm, nota.AtualizadoEm,
	)
	if err != nil {
		return err
	}

	for _, item := range nota.Itens {
		_, err = tx.Exec(ctx, `
			INSERT INTO itens_nota (
				id,
				nota_fiscal_id,
				produto_id,
				produto_codigo,
				produto_descricao,
				quantidade
			)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			item.ID,
			item.NotaID,
			item.ProdutoID,
			item.ProdutoCodigo,
			item.ProdutoDescricao,
			item.Quantidade,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *NotaFiscalRepository) ObterPorID(ctx context.Context, id uuid.UUID) (*entities.NotaFiscal, error) {
	nota := &entities.NotaFiscal{}

	err := r.db.QueryRow(ctx, `
		SELECT id, numero, status, criado_em, atualizado_em,
		       motivo_cancelamento, produto_falhou_id
		FROM notas_fiscais
		WHERE id = $1`, id,
	).Scan(
		&nota.ID,
		&nota.Numero,
		&nota.Status,
		&nota.CriadoEm,
		&nota.AtualizadoEm,
		&nota.MotivoCancelamento,
		&nota.ProdutoFalhouID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerrors.ErrNotaNaoEncontrada
		}
		return nil, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT
			id,
			nota_fiscal_id,
			produto_id,
			produto_codigo,
			produto_descricao,
			quantidade
		FROM itens_nota
		WHERE nota_fiscal_id = $1`, id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item entities.ItemNota
		if err := rows.Scan(
			&item.ID,
			&item.NotaID,
			&item.ProdutoID,
			&item.ProdutoCodigo,
			&item.ProdutoDescricao,
			&item.Quantidade,
		); err != nil {
			return nil, err
		}
		nota.Itens = append(nota.Itens, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nota, nil
}

func (r *NotaFiscalRepository) Listar(ctx context.Context) ([]*entities.NotaFiscal, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, numero, status, criado_em, atualizado_em,
		       motivo_cancelamento, produto_falhou_id
		FROM notas_fiscais
		ORDER BY numero DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notas []*entities.NotaFiscal

	for rows.Next() {
		n := &entities.NotaFiscal{}
		if err := rows.Scan(
			&n.ID,
			&n.Numero,
			&n.Status,
			&n.CriadoEm,
			&n.AtualizadoEm,
			&n.MotivoCancelamento,
			&n.ProdutoFalhouID,
		); err != nil {
			return nil, err
		}

		itensRows, err := r.db.Query(ctx, `
			SELECT
				id,
				nota_fiscal_id,
				produto_id,
				produto_codigo,
				produto_descricao,
				quantidade
			FROM itens_nota
			WHERE nota_fiscal_id = $1`, n.ID,
		)
		if err != nil {
			return nil, err
		}

		for itensRows.Next() {
			var item entities.ItemNota
			if err := itensRows.Scan(
				&item.ID,
				&item.NotaID,
				&item.ProdutoID,
				&item.ProdutoCodigo,
				&item.ProdutoDescricao,
				&item.Quantidade,
			); err != nil {
				itensRows.Close()
				return nil, err
			}
			n.Itens = append(n.Itens, item)
		}

		if err := itensRows.Err(); err != nil {
			itensRows.Close()
			return nil, err
		}

		itensRows.Close()
		notas = append(notas, n)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notas, nil
}

func (r *NotaFiscalRepository) ListarPaginado(
	ctx context.Context,
	filtro ports.ListarNotasFiltro,
) ([]*entities.NotaFiscal, int64, error) {

	offset := (filtro.Page - 1) * filtro.PageSize

	// 🔒 whitelist para evitar SQL injection
	sortByMap := map[string]string{
		"numero":        "nf.numero",
		"status":        "nf.status",
		"criado_em":     "nf.criado_em",
		"atualizado_em": "nf.atualizado_em",
	}

	sortBy := sortByMap[filtro.SortBy]
	if sortBy == "" {
		sortBy = "nf.criado_em"
	}

	sortDir := "DESC"
	if strings.EqualFold(filtro.SortDirection, "asc") {
		sortDir = "ASC"
	}

	search := strings.TrimSpace(filtro.Search)

	where := ""
	args := []any{}
	argPos := 1

	// 🔍 filtro de busca
	if search != "" {
		where = fmt.Sprintf(`
			WHERE (
				CAST(nf.numero AS TEXT) ILIKE $%d
				OR nf.status ILIKE $%d
				OR EXISTS (
					SELECT 1
					FROM itens_nota i
					WHERE i.nota_fiscal_id = nf.id
					  AND (
						i.produto_codigo ILIKE $%d
						OR i.produto_descricao ILIKE $%d
					  )
				)
			)
		`, argPos, argPos, argPos, argPos)

		args = append(args, "%"+search+"%")
		argPos++
	}

	// 📊 TOTAL (para paginação)
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM notas_fiscais nf
		%s
	`, where)

	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// 📦 QUERY principal
	query := fmt.Sprintf(`
		SELECT
			nf.id,
			nf.numero,
			nf.status,
			nf.criado_em,
			nf.atualizado_em,
			nf.motivo_cancelamento,
			nf.produto_falhou_id
		FROM notas_fiscais nf
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, where, sortBy, sortDir, argPos, argPos+1)

	args = append(args, filtro.PageSize, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var notas []*entities.NotaFiscal

	for rows.Next() {
		n := &entities.NotaFiscal{}

		if err := rows.Scan(
			&n.ID,
			&n.Numero,
			&n.Status,
			&n.CriadoEm,
			&n.AtualizadoEm,
			&n.MotivoCancelamento,
			&n.ProdutoFalhouID,
		); err != nil {
			return nil, 0, err
		}

		// 🔁 carregar itens (igual ao seu método atual)
		itensRows, err := r.db.Query(ctx, `
			SELECT
				id,
				nota_fiscal_id,
				produto_id,
				produto_codigo,
				produto_descricao,
				quantidade
			FROM itens_nota
			WHERE nota_fiscal_id = $1
		`, n.ID)
		if err != nil {
			return nil, 0, err
		}

		for itensRows.Next() {
			var item entities.ItemNota
			if err := itensRows.Scan(
				&item.ID,
				&item.NotaID,
				&item.ProdutoID,
				&item.ProdutoCodigo,
				&item.ProdutoDescricao,
				&item.Quantidade,
			); err != nil {
				itensRows.Close()
				return nil, 0, err
			}
			n.Itens = append(n.Itens, item)
		}

		if err := itensRows.Err(); err != nil {
			itensRows.Close()
			return nil, 0, err
		}

		itensRows.Close()
		notas = append(notas, n)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return notas, total, nil
}

func (r *NotaFiscalRepository) Atualizar(ctx context.Context, nota *entities.NotaFiscal) error {
	_, err := r.db.Exec(ctx, `
		UPDATE notas_fiscais
		SET status = $1,
		    atualizado_em = $2,
		    motivo_cancelamento = $3,
		    produto_falhou_id = $4
		WHERE id = $5`,
		nota.Status,
		nota.AtualizadoEm,
		nota.MotivoCancelamento,
		nota.ProdutoFalhouID,
		nota.ID,
	)
	return err
}

// TxNotaFiscalRepository — versão que usa pgx.Tx em vez de pool
type TxNotaFiscalRepository struct {
	tx pgx.Tx
}

func NewTxNotaFiscalRepository(tx pgx.Tx) *TxNotaFiscalRepository {
	return &TxNotaFiscalRepository{tx: tx}
}

func (r *TxNotaFiscalRepository) ProximoNumero(ctx context.Context) (int64, error) {
	var numero int64
	err := r.tx.QueryRow(ctx, "SELECT nextval('nota_fiscal_numero_seq')").Scan(&numero)
	return numero, err
}

func (r *TxNotaFiscalRepository) Salvar(ctx context.Context, nota *entities.NotaFiscal) error {
	_, err := r.tx.Exec(ctx, `
		INSERT INTO notas_fiscais (id, numero, status, criado_em, atualizado_em)
		VALUES ($1, $2, $3, $4, $5)`,
		nota.ID, nota.Numero, nota.Status, nota.CriadoEm, nota.AtualizadoEm,
	)
	if err != nil {
		return err
	}

	for _, item := range nota.Itens {
		_, err = r.tx.Exec(ctx, `
			INSERT INTO itens_nota (
				id,
				nota_fiscal_id,
				produto_id,
				produto_codigo,
				produto_descricao,
				quantidade
			)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			item.ID,
			item.NotaID,
			item.ProdutoID,
			item.ProdutoCodigo,
			item.ProdutoDescricao,
			item.Quantidade,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *TxNotaFiscalRepository) ObterPorID(ctx context.Context, id uuid.UUID) (*entities.NotaFiscal, error) {
	nota := &entities.NotaFiscal{}

	err := r.tx.QueryRow(ctx, `
		SELECT id, numero, status, criado_em, atualizado_em,
		       motivo_cancelamento, produto_falhou_id
		FROM notas_fiscais
		WHERE id = $1`, id,
	).Scan(
		&nota.ID,
		&nota.Numero,
		&nota.Status,
		&nota.CriadoEm,
		&nota.AtualizadoEm,
		&nota.MotivoCancelamento,
		&nota.ProdutoFalhouID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerrors.ErrNotaNaoEncontrada
		}
		return nil, err
	}

	rows, err := r.tx.Query(ctx, `
		SELECT
			id,
			nota_fiscal_id,
			produto_id,
			produto_codigo,
			produto_descricao,
			quantidade
		FROM itens_nota
		WHERE nota_fiscal_id = $1`, id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item entities.ItemNota
		if err := rows.Scan(
			&item.ID,
			&item.NotaID,
			&item.ProdutoID,
			&item.ProdutoCodigo,
			&item.ProdutoDescricao,
			&item.Quantidade,
		); err != nil {
			return nil, err
		}
		nota.Itens = append(nota.Itens, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nota, nil
}

func (r *TxNotaFiscalRepository) Listar(ctx context.Context) ([]*entities.NotaFiscal, error) {
	return nil, nil // não usado em transações
}

func (r *TxNotaFiscalRepository) ListarPaginado(
	ctx context.Context,
	filtro ports.ListarNotasFiltro,
) ([]*entities.NotaFiscal, int64, error) {

	offset := (filtro.Page - 1) * filtro.PageSize

	sortByMap := map[string]string{
		"numero":        "nf.numero",
		"status":        "nf.status",
		"criado_em":     "nf.criado_em",
		"atualizado_em": "nf.atualizado_em",
	}

	sortBy := sortByMap[filtro.SortBy]
	if sortBy == "" {
		sortBy = "nf.criado_em"
	}

	sortDir := "DESC"
	if strings.EqualFold(filtro.SortDirection, "asc") {
		sortDir = "ASC"
	}

	search := strings.TrimSpace(filtro.Search)

	where := ""
	args := []any{}
	argPos := 1

	if search != "" {
		where = fmt.Sprintf(`
			WHERE (
				CAST(nf.numero AS TEXT) ILIKE $%d
				OR nf.status ILIKE $%d
				OR EXISTS (
					SELECT 1
					FROM itens_nota i
					WHERE i.nota_fiscal_id = nf.id
					AND (
						i.produto_codigo ILIKE $%d
						OR i.produto_descricao ILIKE $%d
					)
				)
			)
		`, argPos, argPos, argPos, argPos)

		args = append(args, "%"+search+"%")
		argPos++
	}

	// total
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM notas_fiscais nf
		%s
	`, where)

	var total int64
	if err := r.tx.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// dados
	query := fmt.Sprintf(`
		SELECT
			nf.id,
			nf.numero,
			nf.status,
			nf.criado_em,
			nf.atualizado_em,
			nf.motivo_cancelamento,
			nf.produto_falhou_id
		FROM notas_fiscais nf
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, where, sortBy, sortDir, argPos, argPos+1)

	args = append(args, filtro.PageSize, offset)

	rows, err := r.tx.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var notas []*entities.NotaFiscal

	for rows.Next() {
		n := &entities.NotaFiscal{}
		if err := rows.Scan(
			&n.ID,
			&n.Numero,
			&n.Status,
			&n.CriadoEm,
			&n.AtualizadoEm,
			&n.MotivoCancelamento,
			&n.ProdutoFalhouID,
		); err != nil {
			return nil, 0, err
		}

		notas = append(notas, n)
	}

	return notas, total, nil
}

func (r *TxNotaFiscalRepository) Atualizar(ctx context.Context, nota *entities.NotaFiscal) error {
	_, err := r.tx.Exec(ctx, `
		UPDATE notas_fiscais
		SET status = $1,
		    atualizado_em = $2,
		    motivo_cancelamento = $3,
		    produto_falhou_id = $4
		WHERE id = $5`,
		nota.Status,
		nota.AtualizadoEm,
		nota.MotivoCancelamento,
		nota.ProdutoFalhouID,
		nota.ID,
	)
	return err
}

var _ interface {
	ProximoNumero(context.Context) (int64, error)
	Salvar(context.Context, *entities.NotaFiscal) error
	ObterPorID(context.Context, uuid.UUID) (*entities.NotaFiscal, error)
	Listar(context.Context) ([]*entities.NotaFiscal, error)
	Atualizar(context.Context, *entities.NotaFiscal) error
} = (*NotaFiscalRepository)(nil)

var _ interface {
	ProximoNumero(context.Context) (int64, error)
	Salvar(context.Context, *entities.NotaFiscal) error
	ObterPorID(context.Context, uuid.UUID) (*entities.NotaFiscal, error)
	Listar(context.Context) ([]*entities.NotaFiscal, error)
	Atualizar(context.Context, *entities.NotaFiscal) error
} = (*TxNotaFiscalRepository)(nil)