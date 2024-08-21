package contollers

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/logger"
	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/telegram"
)

func UserMw(repo Repository) telegram.Middleware {
	log := logger.NewStructLogger("user-middleware")
	return func(ctx context.Context, update *telegram.Update) (context.Context, error) {
		if update.Update.Message.Command() == "start" {
			return ctx, nil
		}
		user, err := repo.GetUserById(ctx, update.ChatId)
		if err != nil {
			log.Error("fail get user", zap.Error(err))
			return ctx, errors.New("unknown user")
		}

		ctx = UserInCtx(ctx, user)
		return ctx, nil
	}
}
