-- +goose Up
CREATE SEQUENCE IF NOT EXISTS nota_fiscal_numero_seq START 1;

CREATE TABLE notas_fiscais (
    id                   UUID PRIMARY KEY,
    numero               BIGINT NOT NULL UNIQUE DEFAULT nextval('nota_fiscal_numero_seq'),
    status               TEXT NOT NULL DEFAULT 'Aberta',
    criado_em            TIMESTAMPTZ NOT NULL,
    atualizado_em        TIMESTAMPTZ NOT NULL,
    motivo_cancelamento  TEXT,
    produto_falhou_id    UUID
);

CREATE TABLE itens_nota (
    id             UUID PRIMARY KEY,
    nota_fiscal_id UUID NOT NULL REFERENCES notas_fiscais(id) ON DELETE CASCADE,
    produto_id     UUID NOT NULL,
    produto_codigo TEXT NOT NULL,
    produto_descricao TEXT NOT NULL,
    quantidade     NUMERIC(18,4) NOT NULL CHECK (quantidade > 0)
);

CREATE TABLE outbox_events (
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

CREATE TABLE idempotency_keys (
    chave      TEXT PRIMARY KEY,
    resposta   BYTEA NOT NULL,
    criado_em  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Índices
CREATE INDEX idx_notas_status      ON notas_fiscais(status);
CREATE INDEX idx_itens_nota_id     ON itens_nota(nota_fiscal_id);
CREATE INDEX idx_outbox_status     ON outbox_events(status, proxima_tentativa)
    WHERE status = 'pendente';
CREATE INDEX idx_idempotency_chave ON idempotency_keys(chave);

-- +goose Down
DROP TABLE IF EXISTS idempotency_keys;
DROP TABLE IF EXISTS outbox_events;
DROP TABLE IF EXISTS itens_nota;
DROP TABLE IF EXISTS notas_fiscais;
DROP SEQUENCE IF EXISTS nota_fiscal_numero_seq;
