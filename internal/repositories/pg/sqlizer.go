package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Sqlizer interface {
	// ToSql returns the SQL query as a string, arguments as an array of interfaces, and an error if any.
	//
	// Returns a string representing the SQL query, an array of interfaces for the arguments, and an error.
	ToSql() (string, []interface{}, error)
}

type SqlizerExecutor struct {
	executor Executor
}

// NewSqlizerExecutor creates a new instance of SqlizerExecutor.
//
// It takes an Executor as a parameter and returns a pointer to a SqlizerExecutor.
func NewSqlizerExecutor(executor Executor) *SqlizerExecutor {
	return &SqlizerExecutor{
		executor: executor,
	}
}

// Exec executes the SQL statement represented by the given Sqlizer and returns
// the command tag and any error encountered.
//
// ctx: The context.Context to use for the execution.
// builder: The Sqlizer representing the SQL statement to execute.
// Returns the pgconn.CommandTag representing the result of the execution and
// any error encountered.
func (d *SqlizerExecutor) Exec(ctx context.Context, builder Sqlizer) (pgconn.CommandTag, error) {
	sql, arguments, builderErr := builder.ToSql()
	if builderErr != nil {
		return pgconn.CommandTag{}, sqlBuildErr(builderErr)
	}

	tag, err := d.executor.Exec(ctx, sql, arguments...)
	if err != nil {
		return pgconn.CommandTag{}, sqlExecErr(err)
	}

	return tag, nil
}

// Query executes the SQL statement represented by the given Sqlizer and returns
// the rows and any error encountered.
//
// ctx: The context.Context to use for the execution.
// builder: The Sqlizer representing the SQL statement to execute.
// Returns pgx.Rows representing the result and any error encountered.
func (d *SqlizerExecutor) Query(ctx context.Context, builder Sqlizer) (pgx.Rows, error) { //nolint:ireturn // cause lib
	sql, arguments, builderErr := builder.ToSql()
	if builderErr != nil {
		return nil, sqlBuildErr(builderErr)
	}

	rows, err := d.executor.Query(ctx, sql, arguments...)
	if err != nil {
		return nil, sqlExecErr(err)
	}

	return rows, nil
}

// QueryRow executes the SQL statement represented by the given Sqlizer and returns
// a pgx.Row representing the result.
//
// ctx: The context.Context to use for the execution.
// builder: The Sqlizer representing the SQL statement to execute.
// Returns a pgx.Row representing the result and any error encountered.
func (d *SqlizerExecutor) QueryRow(ctx context.Context, builder Sqlizer) pgx.Row { //nolint:ireturn // cause lib
	sql, args, builderErr := builder.ToSql()
	if builderErr != nil {
		return &erroredRow{builderErr: sqlBuildErr(builderErr)}
	}

	return d.executor.QueryRow(ctx, sql, args...)
}
