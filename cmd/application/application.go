package application

import (
	"fmt"
	"github.com/Yaroher2442/FamilySyncHub/internal/controllers/helpers"

	"github.com/Yaroher2442/FamilySyncHub/cmd/application/config"
	"github.com/Yaroher2442/FamilySyncHub/internal/controllers/commands"
	"github.com/Yaroher2442/FamilySyncHub/internal/controllers/common"
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
		common.START.Route(commands.NewStartController(txManager, repo)),
		common.MyFamilies.Route(helpers.WithUserInCtxHandler(repo, commands.NewMyFamiliesController(txManager, repo))),             //
		common.NewFamily.Route(helpers.WithUserInCtxHandler(repo, commands.NewCreateFamilyController(txManager, repo))),            //
		common.ChoseFamily.Route(helpers.WithUserInCtxHandler(repo, commands.NewChoseFamilyController(txManager, repo))),           //
		common.ChoseFamily.CallbackRoute(helpers.WithUserInCtxHandler(repo, commands.NewChoseFamilyMenuCallback(txManager, repo))), //
		common.AddInFamily.Route(helpers.WithUserInCtxHandler(repo, commands.NewAddInFamilyController(txManager, repo))),           //
		common.NewCategory.Route(helpers.WithUserInCtxHandler(repo, commands.NewCreateCategoryController(txManager, repo))),        //
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
