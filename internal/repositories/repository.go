package repositories

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/Yaroher2442/FamilySyncHub/internal/domain"
	"github.com/Yaroher2442/FamilySyncHub/internal/repositories/pg"
	"github.com/Yaroher2442/FamilySyncHub/internal/repositories/pg/cast"
	"github.com/Yaroher2442/FamilySyncHub/internal/repositories/pg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"strings"
)

func sq() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}

type Repository struct {
	executor *pg.SqlizerExecutor
}

func NewRepository(executor pg.Executor) *Repository {
	return &Repository{
		executor: pg.NewSqlizerExecutor(executor),
	}
}

const (
	USER_TABLE        = "tg_user"
	FAMILY_TABLE      = "family"
	FAMILY_USER_TABLE = "family_user"
)

func ModelUserToDomainUser(model *models.TgUser) *domain.User {
	return &domain.User{
		TgID:           model.TgID,
		ChosenFamilyID: cast.PgUUIDToUUIDPtr(model.ChosenFamilyID),
		FullName:       model.FullName,
		AccountName:    model.AccountName,
	}
}

func (r *Repository) GetUserById(ctx context.Context, id int64) (*domain.User, error) {
	sql := sq().Select("*").
		From(USER_TABLE).
		Where(squirrel.Eq{"tg_id": id}).
		Limit(1)
	model, err := pg.Scan[models.TgUser]().Single(r.executor.Query(ctx, sql))
	if err != nil {
		return nil, err
	}

	return ModelUserToDomainUser(model), err
}

func (r *Repository) GetUserByTgName(ctx context.Context, name string) (*domain.User, error) {
	sql := sq().Select("*").
		From(USER_TABLE).
		Where(squirrel.Eq{"tg_name": name}).
		Limit(1)
	model, err := pg.Scan[models.TgUser]().Single(r.executor.Query(ctx, sql))
	if err != nil {
		return nil, err
	}

	return ModelUserToDomainUser(model), err
}

func (r *Repository) UpdateUser(ctx context.Context, user *domain.User) error {
	_, err := r.executor.Exec(ctx,
		sq().Update(USER_TABLE).
			SetMap(map[string]interface{}{
				"tg_id":            cast.Int64ToPgInt8(user.TgID),
				"account_name":     cast.StrToPgText(user.AccountName),
				"chosen_family_id": cast.UUIDPtrToPgUUID(user.ChosenFamilyID),
				"full_name":        cast.StrToPgText(user.FullName),
			}).Where(squirrel.Eq{"tg_id": user.TgID}))
	return err
}

func (r *Repository) CreateUser(ctx context.Context, user *domain.User) error {
	_, err := r.executor.Exec(ctx,
		sq().Insert(USER_TABLE).
			Columns(
				"tg_id",
				"chosen_family_id",
				"full_name",
				"account_name",
			).
			Values(
				cast.Int64ToPgInt8(user.TgID),
				cast.UUIDPtrToPgUUID(user.ChosenFamilyID),
				cast.StrToPgText(user.FullName),
				cast.StrToPgText(user.AccountName),
			))
	if err != nil && strings.Contains(err.Error(), "violates unique constraint") {
		return domain.ErrUserExists
	}

	return err
}

func (r *Repository) CreateFamily(ctx context.Context, family *domain.Family) error {
	_, err := r.executor.Exec(ctx,
		sq().Insert(FAMILY_TABLE).
			Columns("id", "name").
			Values(
				cast.UUIDToPgUUID(family.ID),
				cast.StrToPgText(family.Name),
			))
	return err
}

func (r *Repository) AddFamilyMember(ctx context.Context, userId domain.USERID, familyId uuid.UUID) error {
	_, err := r.executor.Exec(ctx,
		sq().Insert(FAMILY_USER_TABLE).
			Columns("user_id", "family_id").
			Values(
				cast.Int64ToPgInt8(userId),
				cast.UUIDToPgUUID(familyId),
			))
	return err
}

func (r *Repository) GetFamilyByName(ctx context.Context, name string) (*domain.Family, error) {
	sql := sq().Select("*").
		From(FAMILY_TABLE).
		Where(squirrel.Eq{"name": name}).
		Limit(1)
	model, err := pg.Scan[models.Family]().Single(r.executor.Query(ctx, sql))
	if err != nil {
		return nil, err
	}

	return &domain.Family{
		ID:   cast.PgUUIDToUUID(model.ID),
		Name: model.Name,
	}, err
}

type joinStruct struct {
	UserID   pgtype.Int8 `db:"user_id"`
	FamilyID pgtype.UUID `db:"family_id"`
	ID       pgtype.UUID `db:"id"`
	Name     string      `db:"name"`
}

func (r *Repository) ListUserFamilies(ctx context.Context, user *domain.User) ([]*domain.Family, error) {
	sql := sq().
		Select("*").
		From(FAMILY_USER_TABLE).
		Where(squirrel.Eq{"user_id": user.TgID}).
		Join(FAMILY_TABLE + " ON " + FAMILY_USER_TABLE + ".family_id = " + FAMILY_TABLE + ".id")
	rows, err := r.executor.Query(ctx, sql)
	resultModels, err := pg.Scan[joinStruct]().Multi(rows, err)
	if err != nil {
		return nil, err
	}

	families := make([]*domain.Family, 0, len(resultModels))
	for _, model := range resultModels {
		families = append(families, &domain.Family{
			ID:   cast.PgUUIDToUUID(model.ID),
			Name: model.Name,
		})
	}

	return families, nil
}
