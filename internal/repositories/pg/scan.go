package pg

import (
	"github.com/jackc/pgx/v5"
)

type SingleScanner[T any] func(pgx.Row) (*T, error)

// defaultReflectSingleScanner is a generic function that reflects a single value
// from a pgx.Row into a struct of type T.
//
// It takes a pgx.Row as input and returns a pointer to a struct of type T and an error.
// The returned struct is populated with the value reflected from the row.
// If an error occurs during the reflection process, the function returns nil for the struct and the error.
//
// Parameters:
// - row: A pgx.Row representing the row from which to reflect the value.
//
// Returns:
// - *T: A pointer to a struct of type T populated with the reflected value.
// - error: An error if an error occurs during the reflection process.
func defaultReflectSingleScanner[T any](row pgx.Row) (*T, error) {
	value := new(T)

	err := row.Scan(value)
	if err != nil {
		return nil, sqlScanErr(err)
	}

	return value, nil
}

type MultiScanner[T any] func(pgx.Rows) ([]*T, error)

// defaultReflectMultiScanner is a generic function
// that reflects multiple values from a pgx.Rows into a slice of structs of type T.
//
// It takes a pgx.Rows as input and returns a slice of pointers to structs of type T and an error.
// The returned slice is populated with the values reflected from the rows.
// If an error occurs during the reflection process, the function returns nil for the slice and the error.
//
// Parameters:
// - rows: A pgx.Rows representing the rows from which to reflect the values.
//
// Returns:
// - []*T: A slice of pointers to structs of type T populated with the reflected values.
// - error: An error if an error occurs during the reflection process.
func defaultReflectMultiScanner[T any](rows pgx.Rows) ([]*T, error) {
	values, err := pgx.CollectRows[*T](rows, pgx.RowToStructByName[*T])
	if err != nil {
		return nil, sqlScanErr(err)
	}

	return values, nil
}

type Scanner[T any] struct {
	single SingleScanner[T]
	multi  MultiScanner[T]
}

// Single reflects a single value from a pgx.Row into a struct of type T.
//
// Parameters:
// - row: A pgx.Row representing the row from which to reflect the value.
//
// Returns:
// - *T: A pointer to a struct of type T populated with the reflected value.
// - error: An error if an error occurs during the reflection process.
func (s Scanner[T]) Single(row pgx.Row) (*T, error) {
	return s.single(row)
}

// Multi reflects multiple values from a pgx.Rows into a slice of structs of type T.
//
// Parameters:
// - rows: A pgx.Rows representing the rows from which to reflect the values.
// - err: An error that occurred during the reflection process.
//
// Returns:
// - []*T: A slice of pointers to structs of type T populated with the reflected values.
// - error: An error if an error occurs during the reflection process.
func (s Scanner[T]) Multi(rows pgx.Rows, err error) ([]*T, error) {
	if err != nil {
		return nil, sqlScanErr(err)
	}

	return s.multi(rows)
}

type ScannerOpts[T any] func(*Scanner[T])

// WithSingleSanner is a function that takes a SingleScanner[T] as input and returns a ScannerOpts[T].
//
// It takes a SingleScanner[T] as input and returns a function that takes a *Scanner[T] as
// input and modifies its single field to be the input SingleScanner[T].
//
// Parameters:
// - single: A SingleScanner[T] representing the scanner to be set.
//
// Returns:
// - ScannerOpts[T]: A function that takes a *Scanner[T] as
// input and modifies its single field to be the input SingleScanner[T].
func WithSingleSanner[T any](single SingleScanner[T]) ScannerOpts[T] {
	return func(s *Scanner[T]) {
		s.single = single
	}
}

// WithMultiScanner is a function that takes a MultiScanner[T] as input and returns a ScannerOpts[T].
//
// It takes a MultiScanner[T] as input and returns a function that takes a *Scanner[T] as
// input and modifies its multi field to be the input MultiScanner[T].
//
// Parameters:
// - multi: A MultiScanner[T] representing the scanner to be set.
//
// Returns:
// - ScannerOpts[T]: A function that takes a *Scanner[T] as
// input and modifies its multi field to be the input MultiScanner[T].
func WithMultiScanner[T any](multi MultiScanner[T]) ScannerOpts[T] {
	return func(s *Scanner[T]) {
		s.multi = multi
	}
}

// Scan is a function that creates a Scanner[T] based on the provided ScannerOpts[T] options.
//
// Parameters:
// - opts: Variadic ScannerOpts[T] options.
// Returns:
// - *Scanner[T]: A pointer to the created Scanner[T].
func Scan[T any](opts ...ScannerOpts[T]) *Scanner[T] {
	scanner := &Scanner[T]{
		single: defaultReflectSingleScanner[T],
		multi:  defaultReflectMultiScanner[T],
	}

	for _, opt := range opts {
		opt(scanner)
	}

	return scanner
}
