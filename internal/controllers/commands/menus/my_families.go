package menus

import (
	"github.com/Yaroher2442/FamilySyncHub/internal/controllers/common"
	"github.com/Yaroher2442/FamilySyncHub/internal/domain"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ChoseFamilyMenu(families []*domain.Family) tgbotapi.InlineKeyboardMarkup {
	btns := make([]tgbotapi.InlineKeyboardButton, 0, len(families))
	for _, family := range families {
		btns = append(btns, common.ChoseFamily.InlineKeyboardButtonWithText(
			family.Name,
			common.ChoseFamily.CallbackPayload(family.ID.String()),
		))
	}

	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(btns...))
}
