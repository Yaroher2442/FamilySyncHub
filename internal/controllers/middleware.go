package controllers

import (
	"context"
	"errors"
	"github.com/Yaroher2442/FamilySyncHub/internal/controllers/commands"
	"github.com/Yaroher2442/FamilySyncHub/internal/controllers/helpers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"go.uber.org/zap"

	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/logger"
	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/telegram"
)

func UserMw(repo commands.Repository) telegram.Middleware {
	log := logger.NewStructLogger("user-middleware")
	return func(ctx context.Context, update *telegram.Update) (context.Context, error) {
		if update.Update.Message.Command() == "start" {
			return ctx, nil
		}
		user, err := repo.GetUserById(ctx, update.ChatId)
		if err != nil {
			sentErr := helpers.OnlyErr(
				update.Bot.Send(
					tgbotapi.NewMessage(
						update.ChatId,
						"you are not registered yet please use /start command",
					),
				),
			)
			if sentErr != nil {
				log.Error("fail send auth message message", zap.Error(sentErr))
				return ctx, sentErr
			}
			log.Error("fail get user", zap.Error(err))
			return ctx, errors.New("unknown user")
		}

		return helpers.UserInCtx(ctx, user), nil
	}
}
