package contollers

import (
	"context"
	"fmt"
	"github.com/Yaroher2442/FamilySyncHub/internal/domain"
)

type UserContextKey struct{}

var ErrNoUser = fmt.Errorf("no user in context")

func UserFromCtx(ctx context.Context) (*domain.User, error) {
	user, ok := ctx.Value(UserContextKey{}).(*domain.User)
	if !ok {
		return nil, ErrNoUser
	}
	return user, nil
}

func UserInCtx(ctx context.Context, user *domain.User) context.Context {
	return context.WithValue(ctx, UserContextKey{}, user)
}
