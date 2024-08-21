package cli

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Yaroher2442/FamilySyncHub/cmd/application"
	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/shutdown"
)

func makeQuitSignal() chan os.Signal {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	return quit
}

func run() error {
	app, err := application.New()
	if err != nil {
		return fmt.Errorf("fail init app: %w", err)
	}

	if err := app.Run(); err != nil {
		return fmt.Errorf("fail run app: %w", err)
	}

	shutdown.WaitSignal(makeQuitSignal())

	return nil
}
