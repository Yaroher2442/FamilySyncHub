package pg

import (
	"errors"
	"fmt"
)

const (
	sqlBuildStageErr = "sql build stage: "
	sqlExecStageErr  = "sql exec stage: "
	sqlScanStageErr  = "sql scan stage: "
)

type erroredRow struct {
	builderErr error
}

func (b *erroredRow) Scan(...any) error {
	return b.builderErr
}

type executorError struct {
	err   error
	stage string
}

func (e *executorError) Error() string {
	return fmt.Sprintf("%s %s", e.stage, e.err)
}

func sqlBuildErr(err error) error {
	return &executorError{
		err:   err,
		stage: sqlBuildStageErr,
	}
}

func sqlExecErr(err error) error {
	return &executorError{
		err:   err,
		stage: sqlExecStageErr,
	}
}

func sqlScanErr(err error) error {
	return &executorError{
		err:   err,
		stage: sqlScanStageErr,
	}
}

func IsPgxErr(err, target error) bool {
	var cast *executorError
	ok := errors.As(err, &cast)
	if ok {
		return errors.Is(cast.err, target)
	}

	return errors.Is(err, target)
}
