package contollers

import (
	"context"
	"errors"
	"fmt"
	"github.com/Yaroher2442/FamilySyncHub/internal/contollers/menus"
	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/telegram"
	"github.com/avito-tech/go-transaction-manager/trm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"github.com/Yaroher2442/FamilySyncHub/internal/domain"
	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/logger"
)

type StartController struct {
	logger *zap.Logger
	repo   Repository
	txm    trm.Manager
}

func NewStartController(txm trm.Manager, repository Repository) *StartController {
	return &StartController{
		logger: logger.NewStructLogger("start-controller"),
		repo:   repository,
		txm:    txm,
	}
}

func (s *StartController) Handle(ctx context.Context, update *telegram.Update) error {
	if update.Update.FromChat().Type != "private" {
		_, err := update.Bot.Send(tgbotapi.NewMessage(update.ChatId, "Only for private chats"))

		return err
	}
	user := &domain.User{
		TgID:           update.ChatId,
		AccountName:    update.Update.FromChat().UserName,
		ChosenFamilyID: nil,
		FullName: fmt.Sprintf(
			"%s %s",
			update.Update.FromChat().FirstName,
			update.Update.FromChat().LastName,
		),
	}
	text := "Welcome!"
	err := s.repo.CreateUser(ctx, user)
	if err != nil {
		if !errors.Is(err, domain.ErrUserExists) {
			return err
		}
		text = "User already exists"
	}

	msg := tgbotapi.NewMessage(update.ChatId, text)
	msg.ReplyMarkup = menus.StartMenu()
	_, err = update.Bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
