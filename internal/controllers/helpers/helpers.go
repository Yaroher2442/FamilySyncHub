package helpers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Yaroher2442/FamilySyncHub/internal/domain"
	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/iancoleman/strcase"
)

type UserContextKey string

const (
	UserCtxKey UserContextKey = "user"
)

var ErrNoUser = fmt.Errorf("no user in context")

func UserFromCtx(ctx context.Context) (*domain.User, error) {
	user, ok := ctx.Value(UserCtxKey).(*domain.User)
	if !ok {
		return nil, ErrNoUser
	}
	return user, nil
}

func UserInCtx(ctx context.Context, user *domain.User) context.Context {
	return context.WithValue(ctx, UserCtxKey, user)
}

func OnlyErr(_ any, err error) error {
	return err
}

func CamelCaseArg(update *telegram.Update) string {
	return strcase.ToCamel(update.Update.Message.CommandArguments())
}

func IsArgEmpty(update *telegram.Update) bool {
	return update.Update.Message.CommandArguments() == ""
}

type UserGetter interface {
	GetUserById(ctx context.Context, id int64) (*domain.User, error)
}

type ctxUserRoute struct {
	userGetter UserGetter
	telegram.Handler
}

func (c ctxUserRoute) Handle(ctx context.Context, update *telegram.Update) error {
	user, err := c.userGetter.GetUserById(ctx, update.Update.FromChat().ID)
	if errors.Is(err, domain.ErrNoUser) {
		return OnlyErr(update.Bot.Send(tgbotapi.NewMessage(update.Update.FromChat().ID, "Please, use /start to get started")))
	}
	if err != nil {
		return err
	}

	return c.Handler.Handle(UserInCtx(ctx, user), update)
}

func WithUserInCtxHandler(userGetter UserGetter, handler telegram.Handler) telegram.Handler {
	return &ctxUserRoute{userGetter, handler}
}

func StructCallbackPayload[T any](payload T) (string, error) {
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	return string(payloadJson), nil
}

func ParseStructCallbackPayload[T any](payload string) (T, error) {
	var result T
	err := json.Unmarshal([]byte(payload), &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func ParseStructCallbackPayloadWithErr[T any](payload string, err error) (T, error) {
	if err != nil {
		return *new(T), err
	}

	return ParseStructCallbackPayload[T](payload)
}
