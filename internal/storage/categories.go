package storage

import (
	"context"

	"github.com/Onnywrite/nwstep/internal/domain/models"
	"github.com/Onnywrite/nwstep/internal/lib/cuteql"
	"github.com/google/uuid"
)

func (pg *PgStorage) Categories(ctx context.Context) ([]models.Category, error) {
	categories, tx, err := cuteql.Query[models.Category](ctx, pg.db, `SELECT * FROM categories`)
	if err != nil {
		return nil, err
	}

	return categories, cuteql.Commit(tx)
}

func (pg *PgStorage) CategoryById(ctx context.Context, id int) (*models.Category, error) {
	category, tx, err := cuteql.Get[models.Category](ctx, pg.db, `
	SELECT * FROM categories
	WHERE category_id = $1
	`, id)
	if err != nil {
		return nil, err
	}

	return category, cuteql.Commit(tx)
}

func (pg *PgStorage) Rating(ctx context.Context, uid uuid.UUID, categId int) (int, error) {
	rating, tx, err := cuteql.Get[int](ctx, pg.db, `
	SELECT rating FROM ratings
	WHERE user_id = $1 AND category_id = $2
	`, uid, categId)
	if err != nil {
		return 0, err
	}

	return *rating, cuteql.Commit(tx)
}

func (pg *PgStorage) CategoryTop(ctx context.Context, categoryId, limit int) ([]models.TopUser, error) {
	top, tx, err := cuteql.Query[models.TopUser](ctx, pg.db, `
	SELECT
		ROW_NUMBER() OVER (ORDER BY rating DESC) AS position,
		users.user_id AS user_id,
		users.nickname AS nickname,
		ratings.rating AS rating
	FROM ratings
	JOIN users ON ratings.user_id = users.user_id
	WHERE ratings.category_id = $1
	LIMIT $2
	`, categoryId, limit)
	if err != nil {
		return nil, err
	}

	return top, cuteql.Commit(tx)
}

func (pg *PgStorage) UserTopPosition(ctx context.Context, categoryId int, userId uuid.UUID) (*models.TopUser, error) {
	userTop, tx, err := cuteql.Get[models.TopUser](ctx, pg.db, `
	WITH top_users AS (
		SELECT
			ROW_NUMBER() OVER (ORDER BY rating DESC) AS position,
			users.user_id AS user_id,
			users.nickname AS nickname,
			ratings.rating AS rating
		FROM ratings
		JOIN users ON ratings.user_id = users.user_id
		WHERE ratings.category_id = $1
	)
	SELECT * FROM top_users
	WHERE user_id = $2
	`, categoryId, userId)
	if err != nil {
		return nil, err
	}

	return userTop, cuteql.Commit(tx)
}
