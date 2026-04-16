package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nf-system/servico-faturamento/internal/application/ports"
)

type PgUnitOfWork struct {
	db *pgxpool.Pool
}

func NewUnitOfWork(db *pgxpool.Pool) *PgUnitOfWork {
	return &PgUnitOfWork{db: db}
}

func (u *PgUnitOfWork) Execute(ctx context.Context, fn func(repos ports.TxRepositories) error) error {
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	repos := ports.TxRepositories{
		Notas:  NewTxNotaFiscalRepository(tx),
		Outbox: NewTxOutboxRepository(tx),
	}

	if err := fn(repos); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
