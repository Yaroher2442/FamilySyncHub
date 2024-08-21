package application

import (
	"fmt"
	"github.com/Yaroher2442/FamilySyncHub/cmd/application/config"
	"github.com/Yaroher2442/FamilySyncHub/internal/contollers"
	"github.com/Yaroher2442/FamilySyncHub/internal/contollers/commands"
	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/logger"
	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/shutdown"
	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/telegram"
	"github.com/Yaroher2442/FamilySyncHub/internal/repositories"
	"github.com/Yaroher2442/FamilySyncHub/internal/repositories/pg"
)

type Application struct {
	Listener *telegram.Listener
}

func New() (*Application, error) {
	appConfig, err := config.NewLoader[config.Config]().Load()
	if err != nil {
		return nil, fmt.Errorf("err init config: %w", err)
	}
	psql, closePsql, err := pg.NewPsql(appConfig.Postgres)
	if err != nil {
		return nil, fmt.Errorf("err init psql: %w", err)
	}

	shutdown.RegisterFn(func() {
		closePsql()
		logger.Info("Close psql")
	})

	executor, txManager, err := pg.NewTxFlow(psql)
	if err != nil {
		return nil, fmt.Errorf("pgtx.NewTxFlow: %w", err)
	}

	repo := repositories.NewRepository(executor)

	routes := telegram.NewRouter(
		appConfig.Telegram,
		contollers.NewUnknownController(),
		telegram.Middlewares(
			contollers.UserMw(repo),
		),
		commands.START.Route(contollers.NewStartController(txManager, repo)),
		commands.MyFamilies.Route(contollers.NewMyFamiliesController(txManager, repo)),
		commands.ChoseFamily.Route(contollers.NewChoseFamilyController(txManager, repo)),
		commands.AddInFamily.Route(contollers.NewAddInFamilyController(txManager, repo)),
	)

	return &Application{
		Listener: telegram.NewListener(appConfig.Telegram, routes),
	}, nil
}

func (app *Application) Run() error {
	closeFn, err := app.Listener.Start()
	if err != nil {
		return fmt.Errorf("fail start listener: %w", err)
	}
	shutdown.RegisterFn(closeFn)
	return nil
}
