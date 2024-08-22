package controllers

import (
	"context"

	"go.uber.org/zap"

	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/logger"
	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UnknownController struct {
	logger *zap.Logger
}

func NewUnknownController() *UnknownController {
	return &UnknownController{
		logger: logger.NewStructLogger("unknown-controller"),
	}
}

func (u UnknownController) Handle(ctx context.Context, update *telegram.Update) error {
	_, err := update.Bot.Send(tgbotapi.NewMessage(update.ChatId, "Unknown text/command"))
	return err
}
