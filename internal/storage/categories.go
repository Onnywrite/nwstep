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
