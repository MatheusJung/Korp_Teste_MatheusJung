package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nf-system/servico-faturamento/internal/domain/entities"
)

// ── Outbox ────────────────────────────────────────────────────────────────────

type OutboxRepository struct {
	db *pgxpool.Pool
}

func NewOutboxRepository(db *pgxpool.Pool) *OutboxRepository {
	return &OutboxRepository{db: db}
}

func (r *OutboxRepository) Salvar(ctx context.Context, event *entities.OutboxEvent) error {
	payload, err := json.Marshal(event.Payload)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, `
		INSERT INTO outbox_events
		  (id, nota_fiscal_id, payload, status, tentativas, max_tentativas, criado_em, proxima_tentativa)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		event.ID, event.NotaFiscalID, payload, event.Status,
		event.Tentativas, event.MaxTentativas, event.CriadoEm, event.ProximaTentativa,
	)
	return err
}

// BuscarPendentes usa FOR UPDATE SKIP LOCKED para que múltiplos workers
// nunca processem o mesmo evento ao mesmo tempo.
func (r *OutboxRepository) BuscarPendentes(ctx context.Context, limite int) ([]*entities.OutboxEvent, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, nota_fiscal_id, payload, status, tentativas, max_tentativas,
		       criado_em, processado_em, proxima_tentativa
		FROM outbox_events
		WHERE status = 'pendente'
		  AND proxima_tentativa <= NOW()
		ORDER BY proxima_tentativa
		LIMIT $1
		FOR UPDATE SKIP LOCKED`,
		limite,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var eventos []*entities.OutboxEvent
	for rows.Next() {
		e := &entities.OutboxEvent{}
		var payloadBytes []byte
		if err := rows.Scan(
			&e.ID, &e.NotaFiscalID, &payloadBytes, &e.Status,
			&e.Tentativas, &e.MaxTentativas,
			&e.CriadoEm, &e.ProcessadoEm, &e.ProximaTentativa,
		); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(payloadBytes, &e.Payload); err != nil {
			return nil, err
		}
		eventos = append(eventos, e)
	}
	return eventos, nil
}

func (r *OutboxRepository) Atualizar(ctx context.Context, event *entities.OutboxEvent) error {
	payload, err := json.Marshal(event.Payload)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, `
		UPDATE outbox_events
		SET status = $1, tentativas = $2, payload = $3,
		    processado_em = $4, proxima_tentativa = $5
		WHERE id = $6`,
		event.Status, event.Tentativas, payload,
		event.ProcessadoEm, event.ProximaTentativa,
		event.ID,
	)
	return err
}

// TxOutboxRepository — versão transacional para UnitOfWork
type TxOutboxRepository struct {
	tx pgx.Tx
}

func NewTxOutboxRepository(tx pgx.Tx) *TxOutboxRepository {
	return &TxOutboxRepository{tx: tx}
}

func (r *TxOutboxRepository) Salvar(ctx context.Context, event *entities.OutboxEvent) error {
	payload, err := json.Marshal(event.Payload)
	if err != nil {
		return err
	}
	_, err = r.tx.Exec(ctx, `
		INSERT INTO outbox_events
		  (id, nota_fiscal_id, payload, status, tentativas, max_tentativas, criado_em, proxima_tentativa)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		event.ID, event.NotaFiscalID, payload, event.Status,
		event.Tentativas, event.MaxTentativas, event.CriadoEm, event.ProximaTentativa,
	)
	return err
}

func (r *TxOutboxRepository) BuscarPendentes(ctx context.Context, limite int) ([]*entities.OutboxEvent, error) {
	return nil, nil // não usado em transações
}

func (r *TxOutboxRepository) Atualizar(ctx context.Context, event *entities.OutboxEvent) error {
	payload, err := json.Marshal(event.Payload)
	if err != nil {
		return err
	}
	_, err = r.tx.Exec(ctx, `
		UPDATE outbox_events
		SET status = $1, tentativas = $2, payload = $3,
		    processado_em = $4, proxima_tentativa = $5
		WHERE id = $6`,
		event.Status, event.Tentativas, payload,
		event.ProcessadoEm, event.ProximaTentativa,
		event.ID,
	)
	return err
}

// ── Idempotency ───────────────────────────────────────────────────────────────

type IdempotencyRepository struct {
	db *pgxpool.Pool
}

func NewIdempotencyRepository(db *pgxpool.Pool) *IdempotencyRepository {
	return &IdempotencyRepository{db: db}
}

func (r *IdempotencyRepository) Existe(ctx context.Context, chave string) (bool, error) {
	var existe bool
	err := r.db.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM idempotency_keys WHERE chave = $1)", chave,
	).Scan(&existe)
	return existe, err
}

func (r *IdempotencyRepository) Registrar(ctx context.Context, chave string, resposta []byte) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO idempotency_keys (chave, resposta, criado_em)
		VALUES ($1, $2, $3)
		ON CONFLICT (chave) DO NOTHING`,
		chave, resposta, time.Now().UTC(),
	)
	return err
}

func (r *IdempotencyRepository) ObterResposta(ctx context.Context, chave string) ([]byte, error) {
	var resposta []byte
	err := r.db.QueryRow(ctx,
		"SELECT resposta FROM idempotency_keys WHERE chave = $1", chave,
	).Scan(&resposta)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return resposta, err
}

