package pg

import (
	"context"
	"fmt"
	trmpgx "github.com/avito-tech/go-transaction-manager/pgxv5"
	"github.com/avito-tech/go-transaction-manager/trm/manager"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPsql(config *Config) (*pgxpool.Pool, func(), error) {
	ctx := context.Background()

	uri := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		config.Username, config.Password, config.Host, config.Port, config.Database,
	)

	pool, err := pgxpool.New(ctx, uri)
	if err != nil {
		return nil, nil, fmt.Errorf("pgxpool.New: %w", err)
	}

	return pool, func() {
		pool.Close()
	}, nil
}

// NewTXManager creates TX manager with pgxpool.Pool.
// Put constructor for TX manage here since manager depends on pgx5 driver
func NewTXManager(pool *pgxpool.Pool) (*manager.Manager, error) {
	m, err := manager.New(trmpgx.NewDefaultFactory(pool))
	if err != nil {
		return nil, fmt.Errorf("manager.New: %w", err)
	}

	return m, nil
}

// NewTxFlow creates a new transaction flow with the given pgx pool and options.
//
// It takes a pgx pool and a variadic number of TxExecutorOption as arguments.
// It returns a pointer to a TxExecutor, a pointer to a manager.Manager, and an error.
func NewTxFlow(psql *pgxpool.Pool, options ...TxExecutorOption) (*TxExecutor, *manager.Manager, error) {
	txManager, err := NewTXManager(psql)
	if err != nil {
		return nil, nil, fmt.Errorf("postgres.NewTXManager: %w", err)
	}

	return NewTxExecutor(psql, options...), txManager, nil
}
