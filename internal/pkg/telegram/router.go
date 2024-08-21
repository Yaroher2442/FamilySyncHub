package telegram

import (
	"context"
	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Handler interface {
	Handle(ctx context.Context, update *Update) error
}

type Route interface {
	Match(update *tgbotapi.Update) bool
	Handler
}

type UnknownHandler Handler

type StopFlow = bool

type Middleware func(ctx context.Context, update *Update) (context.Context, error)

type Router struct {
	unknown     UnknownHandler
	routes      []Route
	middlewares []Middleware
	logger      *zap.Logger
	addErr      bool
}

func NewRouter(config *Config, unknown UnknownHandler, middlewares []Middleware, routes ...Route) *Router {
	return &Router{
		unknown:     unknown,
		middlewares: middlewares,
		logger:      logger.NewStructLogger("telegram-router"),
		routes:      routes,
		addErr:      config.Debug,
	}
}

func (f *Router) Update(ctx context.Context, update *Update) {
	targetCtx := ctx
	for _, mw := range f.middlewares {
		mwCtx, err := mw(targetCtx, update)
		if err != nil {
			f.logger.Error("fail run middleware", zap.Error(err))
			_, sendErr := update.Bot.Send(
				tgbotapi.NewMessage(update.ChatId, "something went wrong, please try again"),
			)
			if sendErr != nil {
				f.logger.Error("fail send failing message", zap.Error(sendErr))
			}
			return
		}
		targetCtx = mwCtx
	}

	handle := f.unknown

	for _, route := range f.routes {
		if route.Match(update.Update) {
			handle = route
			break
		}
	}

	err := handle.Handle(ctx, update)
	if err != nil {
		f.logger.Error("fail handle message", zap.Error(err))
		text := "something went wrong, please try again"
		if f.addErr {
			text = text + "\n" + err.Error()
		}
		_, sendErr := update.Bot.Send(
			tgbotapi.NewMessage(update.ChatId, text),
		)
		if sendErr != nil {
			f.logger.Error("fail send failing message", zap.Error(sendErr))
		}
	}
}
