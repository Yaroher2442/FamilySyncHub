package commands

import (
	"context"
	"github.com/Yaroher2442/FamilySyncHub/internal/controllers/helpers"
	"github.com/Yaroher2442/FamilySyncHub/internal/domain"
	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/logger"
	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/telegram"
	"github.com/avito-tech/go-transaction-manager/trm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CreateFamilyController struct {
	logger *zap.Logger
	repo   Repository
	txm    trm.Manager
}

func NewCreateFamilyController(tx trm.Manager, repository Repository) *CreateFamilyController {
	return &CreateFamilyController{
		logger: logger.NewStructLogger("new-family-controller"),
		repo:   repository,
		txm:    tx,
	}
}

func (c *CreateFamilyController) Handle(ctx context.Context, update *telegram.Update) error {
	user, err := helpers.UserFromCtx(ctx)
	if err != nil {
		return err
	}
	if helpers.ArgEmpty(update) {
		return helpers.OnlyErr(
			update.Bot.Send(
				tgbotapi.NewMessage(
					update.ChatId,
					"Enter your family name like this: /new_family MyFamily",
				),
			),
		)
	}
	famName := helpers.FamName(update)
	txErr := c.txm.Do(ctx, func(ctx context.Context) error {
		famId := uuid.New()
		family := &domain.Family{
			ID:   famId,
			Name: famName,
		}
		if err := c.repo.CreateFamily(ctx, family); err != nil {
			return err
		}
		if err := c.repo.AddFamilyMember(ctx, user.TgID, family.ID); err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		return txErr
	}
	return helpers.OnlyErr(update.Bot.Send(tgbotapi.NewMessage(
		update.ChatId,
		"Your family was created successfully with name "+famName,
	)))
}
