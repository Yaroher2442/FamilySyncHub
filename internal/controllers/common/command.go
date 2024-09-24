package common

import (
	"fmt"
	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type Command string

const (
	START       Command = "start"
	MyFamilies  Command = "my_families"
	NewFamily   Command = "new_family"
	ChoseFamily Command = "chose_family"
	AddInFamily Command = "add_in_family"
	NewCategory Command = "new_category"
)

func (c Command) String() string {
	return string(c)
}

func (c Command) WithSlash() string {
	return "/" + string(c)
}

func (c Command) KeyboardButton() tgbotapi.KeyboardButton {
	return tgbotapi.NewKeyboardButton(c.WithSlash())
}

func (c Command) KeyboardButtonWithText(text string) tgbotapi.KeyboardButton {
	return tgbotapi.NewKeyboardButton(c.WithSlash() + " " + text)
}

func (c Command) InlineKeyboardButtonWithText(text string, data string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(text, data)
}

func (c Command) CallbackPayload(payload any) string {
	return fmt.Sprintf("%s::%s", c.String(), payload)
}

func (c Command) ExtractCallbackPayload(updateData *telegram.Update) (string, error) {
	if updateData == nil {
		return "", fmt.Errorf("callback command not found")
	}

	if updateData.Update.CallbackQuery == nil {
		return "", fmt.Errorf("callback command not found")
	}

	if updateData.Update.CallbackQuery.Data == "" {
		return "", fmt.Errorf("callback command not found")
	}

	if !strings.HasPrefix(updateData.Update.CallbackQuery.Data, c.String()) {
		return "", fmt.Errorf("callback command not found")
	}

	elems := strings.Split(updateData.Update.CallbackData(), "::")
	if len(elems) != 2 || elems[0] != c.String() || len(elems[1]) == 0 {
		return "", fmt.Errorf("callback command not found")
	}

	return elems[1], nil
}

func (c Command) Route(handler telegram.Handler) telegram.Route {
	return telegram.CommandRoute(string(c), handler)
}

func (c Command) CallbackRoute(handler telegram.Handler) telegram.Route {
	//return telegram.CallbackRoute(regexp.MustCompile(fmt.Sprintf(`^%s$`, c.String())), handler)
	return telegram.CallbackRoute(fmt.Sprintf("%s::.*", c.String()), handler)
}
