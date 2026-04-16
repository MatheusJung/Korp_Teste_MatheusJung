package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
			INSERT INTO itens_nota (id, nota_fiscal_id, produto_id, quantidade)
			VALUES ($1, $2, $3, $4)`,
			item.ID, item.NotaID, item.ProdutoID, item.Quantidade,
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
		FROM notas_fiscais WHERE id = $1`, id,
	).Scan(
		&nota.ID, &nota.Numero, &nota.Status, &nota.CriadoEm, &nota.AtualizadoEm,
		&nota.MotivoCancelamento, &nota.ProdutoFalhouID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerrors.ErrNotaNaoEncontrada
		}
		return nil, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, nota_fiscal_id, produto_id, quantidade
		FROM itens_nota WHERE nota_fiscal_id = $1`, id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item entities.ItemNota
		if err := rows.Scan(&item.ID, &item.NotaID, &item.ProdutoID, &item.Quantidade); err != nil {
			return nil, err
		}
		nota.Itens = append(nota.Itens, item)
	}

	return nota, nil
}

func (r *NotaFiscalRepository) Listar(ctx context.Context) ([]*entities.NotaFiscal, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, numero, status, criado_em, atualizado_em,
		       motivo_cancelamento, produto_falhou_id
		FROM notas_fiscais ORDER BY numero DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notas []*entities.NotaFiscal

	for rows.Next() {
		n := &entities.NotaFiscal{}
		if err := rows.Scan(
			&n.ID, &n.Numero, &n.Status, &n.CriadoEm, &n.AtualizadoEm,
			&n.MotivoCancelamento, &n.ProdutoFalhouID,
		); err != nil {
			return nil, err
		}

		// 🔥 CARREGAR ITENS AQUI
		itensRows, err := r.db.Query(ctx, `
			SELECT id, nota_fiscal_id, produto_id, quantidade
			FROM itens_nota WHERE nota_fiscal_id = $1`, n.ID,
		)
		if err != nil {
			return nil, err
		}

		for itensRows.Next() {
			var item entities.ItemNota
			if err := itensRows.Scan(&item.ID, &item.NotaID, &item.ProdutoID, &item.Quantidade); err != nil {
				itensRows.Close()
				return nil, err
			}
			n.Itens = append(n.Itens, item)
		}
		itensRows.Close()

		notas = append(notas, n)
	}

	return notas, nil
}

func (r *NotaFiscalRepository) Atualizar(ctx context.Context, nota *entities.NotaFiscal) error {
	_, err := r.db.Exec(ctx, `
		UPDATE notas_fiscais
		SET status = $1, atualizado_em = $2,
		    motivo_cancelamento = $3, produto_falhou_id = $4
		WHERE id = $5`,
		nota.Status, nota.AtualizadoEm,
		nota.MotivoCancelamento, nota.ProdutoFalhouID,
		nota.ID,
	)
	return err
}

// TxNotaFiscalRepository — versão que usa pgx.Tx em vez de pool (para UnitOfWork)
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
			INSERT INTO itens_nota (id, nota_fiscal_id, produto_id, quantidade)
			VALUES ($1, $2, $3, $4)`,
			item.ID, item.NotaID, item.ProdutoID, item.Quantidade,
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
		FROM notas_fiscais WHERE id = $1`, id,
	).Scan(
		&nota.ID, &nota.Numero, &nota.Status, &nota.CriadoEm, &nota.AtualizadoEm,
		&nota.MotivoCancelamento, &nota.ProdutoFalhouID,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domainerrors.ErrNotaNaoEncontrada
	}
	return nota, err
}

func (r *TxNotaFiscalRepository) Listar(ctx context.Context) ([]*entities.NotaFiscal, error) {
	return nil, nil // não usado em transações
}

func (r *TxNotaFiscalRepository) Atualizar(ctx context.Context, nota *entities.NotaFiscal) error {
	_, err := r.tx.Exec(ctx, `
		UPDATE notas_fiscais
		SET status = $1, atualizado_em = $2,
		    motivo_cancelamento = $3, produto_falhou_id = $4
		WHERE id = $5`,
		nota.Status, nota.AtualizadoEm,
		nota.MotivoCancelamento, nota.ProdutoFalhouID,
		nota.ID,
	)
	return err
}

// Garantir que ambos implementam a interface
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

// Suprimir warning de enums não utilizado
