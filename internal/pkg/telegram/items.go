package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
)

func Middlewares(mw ...Middleware) []Middleware {
	return mw
}

type textRouterItem struct {
	pattern *regexp.Regexp
	Handler
}

func (t *textRouterItem) Match(update *tgbotapi.Update) bool {
	return t.pattern.MatchString(update.Message.Text)
}

func TextRoute(pattern string, handler Handler) Route {
	return &textRouterItem{pattern: regexp.MustCompile(pattern), Handler: handler}
}

type commandRouterItem struct {
	command string
	Handler
}

func (c *commandRouterItem) Match(update *tgbotapi.Update) bool {
	return update.Message.IsCommand() && update.Message.Command() == c.command
}

func CommandRoute(pattern string, handler Handler) Route {
	return &commandRouterItem{command: pattern, Handler: handler}
}
