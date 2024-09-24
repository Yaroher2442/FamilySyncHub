package commands

import (
	"context"
	"errors"
	"github.com/Yaroher2442/FamilySyncHub/internal/controllers/helpers"
	"github.com/Yaroher2442/FamilySyncHub/internal/domain"
	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/logger"
	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/telegram"
	"github.com/avito-tech/go-transaction-manager/trm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CreateCategoryController struct {
	logger *zap.Logger
	repo   Repository
	txm    trm.Manager
}

func NewCreateCategoryController(tx trm.Manager, repository Repository) *CreateCategoryController {
	return &CreateCategoryController{
		logger: logger.NewStructLogger("new-category-controller"),
		repo:   repository,
		txm:    tx,
	}
}

func (c *CreateCategoryController) Handle(ctx context.Context, update *telegram.Update) error {
	user, err := helpers.UserFromCtx(ctx)
	if err != nil {
		return err
	}
	if user.ChosenFamilyID == nil {
		return helpers.OnlyErr(
			update.Bot.Send(
				tgbotapi.NewMessage(
					update.ChatId,
					"Choose your family first: /my_families",
				),
			),
		)
	}
	if helpers.IsArgEmpty(update) {
		return helpers.OnlyErr(
			update.Bot.Send(
				tgbotapi.NewMessage(
					update.ChatId,
					"Enter your category name like this: /new_category MyCategory",
				),
			),
		)
	}
	catName := helpers.CamelCaseArg(update)
	catId := uuid.New()
	category := &domain.Category{
		ID:       catId,
		Name:     catName,
		FamilyID: *user.ChosenFamilyID,
	}
	err = c.repo.CreateCategory(ctx, category)
	if errors.Is(err, domain.ErrDuplicateCategory) {
		return helpers.OnlyErr(
			update.Bot.Send(
				tgbotapi.NewMessage(
					update.ChatId,
					"Category with name "+catName+" already exists",
				),
			),
		)
	}
	if err != nil {
		return err
	}

	return helpers.OnlyErr(
		update.Bot.Send(
			tgbotapi.NewMessage(
				update.ChatId,
				"Category with name "+catName+" created",
			),
		),
	)
}
