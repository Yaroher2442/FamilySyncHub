package helpers

import (
	"context"
	"fmt"
	"github.com/Yaroher2442/FamilySyncHub/internal/domain"
	"github.com/Yaroher2442/FamilySyncHub/internal/pkg/telegram"
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

func FamName(update *telegram.Update) string {
	return strcase.ToCamel(update.Update.Message.CommandArguments())
}

func ArgEmpty(update *telegram.Update) bool {
	return update.Update.Message.CommandArguments() == ""
}
