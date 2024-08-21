package pg

import "fmt"

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
