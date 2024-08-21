package menus

import (
	"github.com/Yaroher2442/FamilySyncHub/internal/contollers/commands"
	"github.com/Yaroher2442/FamilySyncHub/internal/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MyFamiliesMenu(families []*domain.Family) tgbotapi.ReplyKeyboardMarkup {
	btns := make([]tgbotapi.KeyboardButton, 0)
	for _, family := range families {
		btns = append(btns, commands.ChoseFamily.KeyboardButtonWithText(family.Name))
	}

	return tgbotapi.NewOneTimeReplyKeyboard(tgbotapi.NewKeyboardButtonRow(btns...))
}
