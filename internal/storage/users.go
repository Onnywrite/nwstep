package storage

import (
	"context"
	"fmt"

	"github.com/Onnywrite/nwstep/internal/domain/models"
	"github.com/Onnywrite/nwstep/internal/lib/cuteql"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (pg *PgStorage) TruncateTableUsers(ctx context.Context) error {
	tx, err := cuteql.Execute(ctx, pg.db, `TRUNCATE TABLE users CASCADE`)
	if err != nil {
		return err
	}

	return cuteql.Commit(tx)
}

func (pg *PgStorage) SaveUser(ctx context.Context, user models.User,
) (*models.User, error) {
	result, tx, err := cuteql.GetSquirreled[models.User](ctx, pg.db,
		squirrel.
			Insert("users").
			Columns("user_nickname", "user_login", "user_password_hash", "user_birthday").
			Values(user.Nickname, user.Login, user.PasswordHash, user.Birthday).
			Suffix("RETURNING *").PlaceholderFormat(squirrel.Dollar),
	)
	if err != nil {
		return nil, err
	}

	return result, cuteql.Commit(tx)
}

func (pg *PgStorage) UpdateUser(ctx context.Context,
	userId uuid.UUID,
	newValues map[string]any,
) error {
	if len(newValues) == 0 {
		return fmt.Errorf("no fields to update: %w", cuteql.ErrEmptyResult)
	}

	if _, ok := newValues["user_id"]; ok {
		return fmt.Errorf("user_id must not be changed: %w", cuteql.ErrInternal)
	}

	tx, err := cuteql.ExecuteSquirreled(ctx, pg.db,
		squirrel.
			Update("users").
			SetMap(newValues).
			Where("user_id = ?", userId).PlaceholderFormat(squirrel.Dollar),
	)
	if err != nil {
		return err
	}

	return cuteql.Commit(tx)
}

func (pg *PgStorage) UserByLogin(ctx context.Context, login string) (*models.User, error) {
	return pg.getUserWhere(ctx, squirrel.Eq{"user_login": login})
}

func (pg *PgStorage) UserByNickname(ctx context.Context, nickname string) (*models.User, error) {
	return pg.getUserWhere(ctx, squirrel.Eq{"user_nickname": nickname})
}

func (pg *PgStorage) UserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return pg.getUserWhere(ctx, squirrel.Eq{"user_id": id})
}

func (pg *PgStorage) getUserWhere(ctx context.Context, where squirrel.Sqlizer,
) (*models.User, error) {
	user, tx, err := cuteql.GetSquirreled[models.User](ctx, pg.db,
		squirrel.
			Select("*").
			From("users").
			Where(where).PlaceholderFormat(squirrel.Dollar),
	)
	if err != nil {
		return nil, err
	}

	return user, cuteql.Commit(tx)
}
