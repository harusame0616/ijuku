package txrunner

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrLockTimeout = errors.New("lock timeout")

type PgxTransactionRunner struct {
	pool *pgxpool.Pool
}

func NewPgxTransactionRunner(pool *pgxpool.Pool) PgxTransactionRunner {
	return PgxTransactionRunner{pool: pool}
}

func (r PgxTransactionRunner) RunInTransaction(ctx context.Context, f func(tx pgx.Tx) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := f(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r PgxTransactionRunner) AcquireAdvisoryLock(ctx context.Context, tx pgx.Tx, lockKey int64) error {
	if _, err := tx.Exec(ctx, "SET LOCAL lock_timeout = '5s'"); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, "SELECT pg_advisory_xact_lock($1)", lockKey); err != nil {
		var pgErr *pgconn.PgError
		// 55P03: lock_not_available
		if errors.As(err, &pgErr) && pgErr.Code == "55P03" {
			return ErrLockTimeout
		}
		return err
	}
	return nil
}
