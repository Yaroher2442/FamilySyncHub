package commands

import (
	"context"
	"github.com/Yaroher2442/FamilySyncHub/internal/domain"
	"github.com/google/uuid"
)

type Repository interface {
	GetUserById(ctx context.Context, id int64) (*domain.User, error)
	GetUserByTgName(ctx context.Context, name string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	CreateUser(ctx context.Context, user *domain.User) error

	CreateFamily(ctx context.Context, family *domain.Family) error
	AddFamilyMember(ctx context.Context, userId domain.USERID, familyId uuid.UUID) error
	GetFamilyByID(ctx context.Context, id uuid.UUID) (*domain.Family, error)
	GetFamilyByName(ctx context.Context, name string) (*domain.Family, error)
	ListUserFamilies(ctx context.Context, user *domain.User) ([]*domain.Family, error)

	CreateCategory(ctx context.Context, category *domain.Category) error
}
