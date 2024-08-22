package telegram

import (
	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
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
	matched := update.Message.IsCommand() && update.Message.Command() == c.command
	logger.Debug(
		"match command",
		zap.String("command_test", c.command),
		zap.Bool("is_command", update.Message.IsCommand()),
		zap.String("text", update.Message.CommandWithAt()),
		zap.Bool("matched", matched),
	)
	return matched
}

func CommandRoute(pattern string, handler Handler) Route {
	return &commandRouterItem{command: pattern, Handler: handler}
}
