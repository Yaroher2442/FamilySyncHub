package commands

import (
	"context"
	"errors"
	"github.com/Yaroher2442/FamilySyncHub/internal/controllers/commands/menus"
	"github.com/Yaroher2442/FamilySyncHub/internal/controllers/helpers"
	"github.com/Yaroher2442/FamilySyncHub/internal/domain"
	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/logger"
	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/telegram"
	"github.com/avito-tech/go-transaction-manager/trm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type MyFamiliesController struct {
	logger *zap.Logger
	repo   Repository
	txm    trm.Manager
}

func NewMyFamiliesController(tx trm.Manager, repository Repository) *MyFamiliesController {
	return &MyFamiliesController{
		logger: logger.NewStructLogger("my-families-controller"),
		repo:   repository,
		txm:    tx,
	}
}

func (c *MyFamiliesController) Handle(ctx context.Context, update *telegram.Update) error {
	user, err := helpers.UserFromCtx(ctx)
	if err != nil {
		return err
	}

	families, err := c.repo.ListUserFamilies(ctx, user)
	if err != nil {
		if errors.Is(err, domain.ErrFamiliesEmpty) {
			return helpers.OnlyErr(update.Bot.Send(
				tgbotapi.NewMessage(update.ChatId, "You don't have any families. Please, create one first: /new_family ...")),
			)
		}

		return err
	}

	msg := tgbotapi.NewMessage(update.ChatId, "My families")
	msg.ReplyMarkup = menus.ChoseFamilyMenu(families)

	return helpers.OnlyErr(update.Bot.Send(msg))
}
