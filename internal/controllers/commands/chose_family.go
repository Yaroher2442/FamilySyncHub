package commands

import (
	"context"
	"github.com/Yaroher2442/FamilySyncHub/internal/controllers/common"
	"github.com/Yaroher2442/FamilySyncHub/internal/controllers/helpers"
	"github.com/Yaroher2442/FamilySyncHub/internal/domain"
	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/logger"
	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/telegram"
	"github.com/avito-tech/go-transaction-manager/trm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
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

	if helpers.IsArgEmpty(update) {
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
		family, err := c.repo.GetFamilyByName(ctx, helpers.CamelCaseArg(update))
		if err != nil {
			return err
		}

		user.ChosenFamilyID = &family.ID

		return c.repo.UpdateUser(ctx, user)
	})
}

type ChoseFamilyMenuCallback struct {
	logger *zap.Logger
	repo   Repository
	txm    trm.Manager
}

func NewChoseFamilyMenuCallback(tx trm.Manager, repository Repository) *ChoseFamilyMenuCallback {
	return &ChoseFamilyMenuCallback{
		logger: logger.NewStructLogger("chose-family-menu-callback"),
		repo:   repository,
		txm:    tx,
	}
}

func (c *ChoseFamilyMenuCallback) Handle(ctx context.Context, update *telegram.Update) error {
	user, err := helpers.UserFromCtx(ctx)
	if err != nil {
		return err
	}

	familyID, err := common.ChoseFamily.ExtractCallbackPayload(update)
	if err != nil {
		return err
	}

	var fam *domain.Family = nil

	txErr := c.txm.Do(ctx, func(ctx context.Context) error {
		famID := uuid.MustParse(familyID)
		if user.ChosenFamilyID == nil || famID.String() != user.ChosenFamilyID.String() {
			user.ChosenFamilyID = &famID
			updateErr := c.repo.UpdateUser(ctx, user)
			if updateErr != nil {
				return updateErr
			}
		}
		fam, err = c.repo.GetFamilyByID(ctx, famID)
		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		return txErr
	}

	msg := tgbotapi.NewEditMessageText(
		update.ChatId,
		update.Update.CallbackQuery.Message.MessageID,
		"Chosen family: "+fam.Name,
	)
	return helpers.OnlyErr(update.Bot.Send(msg))
}
