package telegram

import (
	"context"
	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"regexp"
)

func Middlewares(mw ...Middleware) []Middleware {
	return mw
}

type TextRouteMatcher struct {
	pattern *regexp.Regexp
	Handler
}

func (t *TextRouteMatcher) Match(_ context.Context, update *tgbotapi.Update) bool {
	return t.pattern.MatchString(update.Message.Text)
}

func TextRoute(pattern string, handler Handler) *TextRouteMatcher {
	return &TextRouteMatcher{pattern: regexp.MustCompile(pattern), Handler: handler}
}

type CommandRouteMatcher struct {
	command string
	Handler
}

func (c *CommandRouteMatcher) Match(_ context.Context, update *tgbotapi.Update) bool {
	if update.Message == nil {
		return false
	}
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

func CommandRoute(pattern string, handler Handler) *CommandRouteMatcher {
	return &CommandRouteMatcher{command: pattern, Handler: handler}
}

type CallbackRouteMatcher struct {
	pattern *regexp.Regexp
	Handler
}

func (c *CallbackRouteMatcher) Match(_ context.Context, update *tgbotapi.Update) bool {
	if update.CallbackQuery == nil {
		return false
	}
	match := c.pattern.MatchString(update.CallbackData())
	logger.Debug(
		"match callback",
		zap.String("data", update.CallbackData()),
		zap.Bool("matched", match),
	)
	return match
}

func CallbackRoute(pattern string, handler Handler) *CallbackRouteMatcher {
	return &CallbackRouteMatcher{pattern: regexp.MustCompile(pattern), Handler: handler}
}

type FnRouteMatcher struct {
	fn func(ctx context.Context, update *tgbotapi.Update) bool
	Handler
}

func (f *FnRouteMatcher) Match(ctx context.Context, update *tgbotapi.Update) bool {
	return f.fn(ctx, update)
}

func Fn(fn func(ctx context.Context, update *tgbotapi.Update) bool, handler Handler) *FnRouteMatcher {
	return &FnRouteMatcher{fn: fn, Handler: handler}
}
