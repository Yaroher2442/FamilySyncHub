package menus

import (
	"github.com/Yaroher2442/FamilySyncHub/internal/controllers/common"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			common.MyFamilies.KeyboardButton(),
		),
	)
}
