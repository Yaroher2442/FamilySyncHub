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
	Match(ctx context.Context, update *tgbotapi.Update) bool
	Handler
}

type UnknownHandler Handler

type StopFlow = bool

type Middleware func(ctx context.Context, update *Update) (context.Context, error)

type Router struct {
	routes []Route
	logger *zap.Logger
	addErr bool
}

func NewRouter(config *Config, routes ...Route) *Router {
	return &Router{
		routes: routes,
		logger: logger.NewStructLogger("telegram-router"),
		addErr: config.Debug,
	}
}

func (f *Router) Update(ctx context.Context, update *Update) {
	for _, route := range f.routes {
		if route.Match(ctx, update.Update) {
			err := route.Handle(ctx, update)
			if err != nil {
				f.logger.Error("handle error", zap.Error(err))
				if f.addErr {
					_, _ = update.Bot.Send(tgbotapi.NewMessage(update.ChatId, err.Error()))
				}
			}
			return
		}
	}
	_, _ = update.Bot.Send(tgbotapi.NewMessage(update.ChatId, "Unknown input"))
}
