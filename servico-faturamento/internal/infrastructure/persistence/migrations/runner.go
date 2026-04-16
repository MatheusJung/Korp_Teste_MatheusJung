package migrations

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Run aplica todas as migrations em ordem. Idempotente — usa tabela de controle interna.
func Run(ctx context.Context, db *pgxpool.Pool) error {
	if err := criarTabelaControle(ctx, db); err != nil {
		return fmt.Errorf("criar tabela de controle: %w", err)
	}

	for _, m := range migrations {
		aplicada, err := jaAplicada(ctx, db, m.versao)
		if err != nil {
			return err
		}
		if aplicada {
			continue
		}

		if _, err := db.Exec(ctx, m.sql); err != nil {
			return fmt.Errorf("migration %s: %w", m.versao, err)
		}

		if err := registrar(ctx, db, m.versao); err != nil {
			return err
		}
	}

	return nil
}

func criarTabelaControle(ctx context.Context, db *pgxpool.Pool) error {
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			versao     TEXT PRIMARY KEY,
			aplicada_em TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`)
	return err
}

func jaAplicada(ctx context.Context, db *pgxpool.Pool, versao string) (bool, error) {
	var existe bool
	err := db.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE versao = $1)", versao,
	).Scan(&existe)
	return existe, err
}

func registrar(ctx context.Context, db *pgxpool.Pool, versao string) error {
	_, err := db.Exec(ctx,
		"INSERT INTO schema_migrations (versao) VALUES ($1)", versao)
	return err
}

type migration struct {
	versao string
	sql    string
}

var migrations = []migration{
	{
		versao: "001_initial",
		sql: `
			CREATE SEQUENCE IF NOT EXISTS nota_fiscal_numero_seq START 1;

			CREATE TABLE IF NOT EXISTS notas_fiscais (
				id                   UUID PRIMARY KEY,
				numero               BIGINT NOT NULL UNIQUE,
				status               TEXT NOT NULL DEFAULT 'Aberta',
				criado_em            TIMESTAMPTZ NOT NULL,
				atualizado_em        TIMESTAMPTZ NOT NULL,
				motivo_cancelamento  TEXT,
				produto_falhou_id    UUID
			);

			CREATE TABLE IF NOT EXISTS itens_nota (
				id             UUID PRIMARY KEY,
				nota_fiscal_id UUID NOT NULL REFERENCES notas_fiscais(id) ON DELETE CASCADE,
				produto_id     UUID NOT NULL,
				produto_codigo TEXT NOT NULL,
    			produto_descricao TEXT NOT NULL,
				quantidade     NUMERIC(18,4) NOT NULL CHECK (quantidade > 0)
			);

			CREATE TABLE IF NOT EXISTS outbox_events (
				id                UUID PRIMARY KEY,
				nota_fiscal_id    UUID NOT NULL REFERENCES notas_fiscais(id),
				payload           JSONB NOT NULL,
				status            TEXT NOT NULL DEFAULT 'pendente',
				tentativas        INT NOT NULL DEFAULT 0,
				max_tentativas    INT NOT NULL DEFAULT 5,
				criado_em         TIMESTAMPTZ NOT NULL,
				processado_em     TIMESTAMPTZ,
				proxima_tentativa TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);

			CREATE TABLE IF NOT EXISTS idempotency_keys (
				chave      TEXT PRIMARY KEY,
				resposta   BYTEA NOT NULL,
				criado_em  TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);

			CREATE INDEX IF NOT EXISTS idx_notas_status  ON notas_fiscais(status);
			CREATE INDEX IF NOT EXISTS idx_itens_nota_id ON itens_nota(nota_fiscal_id);
			CREATE INDEX IF NOT EXISTS idx_outbox_status ON outbox_events(status, proxima_tentativa)
				WHERE status = 'pendente';
		`,
	},
}
