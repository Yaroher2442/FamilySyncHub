package contollers

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/logger"
	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/telegram"
)

type AddInFamilyController struct {
	logger *zap.Logger
	repo   Repository
	txm    trm.Manager
}

func NewAddInFamilyController(tx trm.Manager, repository Repository) *AddInFamilyController {
	return &AddInFamilyController{
		logger: logger.NewStructLogger("add-in-family-controller"),
		repo:   repository,
		txm:    tx,
	}
}

func (c *AddInFamilyController) Handle(ctx context.Context, update *telegram.Update) error {
	user, err := UserFromCtx(ctx)
	if err != nil {
		return err
	}

	if user.ChosenFamilyID == nil {
		_, sntErr := update.Bot.Send(tgbotapi.NewMessage(
			update.ChatId,
			"Please, choose your family first: /choose_family",
		))

		return sntErr
	}

	args := update.Update.Message.CommandArguments()
	if args == "" {
		_, sntErr := update.Bot.Send(tgbotapi.NewMessage(
			update.ChatId,
			"Enter your family member name like this: /add_family_member @SomeTgName",
		))
		if sntErr != nil {
			return sntErr
		}
		return nil
	}

	return c.txm.Do(ctx, func(ctx context.Context) error {
		targetUser, getUserErr := c.repo.GetUserByTgName(ctx, args)
		if getUserErr != nil {
			return getUserErr
		}

		return c.repo.AddFamilyMember(ctx, targetUser.TgID, *user.ChosenFamilyID)
	})

}
