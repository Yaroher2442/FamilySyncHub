package commands

import (
	"context"
	"github.com/Yaroher2442/FamilySyncHub/internal/controllers/helpers"
	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/logger"
	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/telegram"
	"github.com/avito-tech/go-transaction-manager/trm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type ChoseFamilyController struct {
	logger *zap.Logger
	repo   Repository
	txm    trm.Manager
}

func NewChoseFamilyController(tx trm.Manager, repository Repository) *ChoseFamilyController {
	return &ChoseFamilyController{
		logger: logger.NewStructLogger("chose-family-controller"),
		repo:   repository,
		txm:    tx,
	}
}

func (c *ChoseFamilyController) Handle(ctx context.Context, update *telegram.Update) error {
	user, err := helpers.UserFromCtx(ctx)
	if err != nil {
		return err
	}

	if helpers.ArgEmpty(update) {
		_, sntErr := update.Bot.Send(tgbotapi.NewMessage(
			update.ChatId,
			"Enter your family name like this: /choose_family MyFamily",
		))
		if sntErr != nil {
			return sntErr
		}
		return nil
	}

	return c.txm.Do(ctx, func(ctx context.Context) error {
		family, err := c.repo.GetFamilyByName(ctx, helpers.FamName(update))
		if err != nil {
			return err
		}

		user.ChosenFamilyID = &family.ID

		return c.repo.UpdateUser(ctx, user)
	})
}
