package commands

import (
	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Command string

const (
	START       Command = "start"
	MyFamilies  Command = "my_families"
	ChoseFamily Command = "chose_family"
	AddInFamily Command = "add_in_family"
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

func (c Command) Route(handler telegram.Handler) telegram.Route {
	return telegram.CommandRoute(string(c), handler)
}
