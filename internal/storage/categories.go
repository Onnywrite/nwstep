package storage

import (
	"context"

	"github.com/Onnywrite/nwstep/internal/domain/models"
	"github.com/Onnywrite/nwstep/internal/lib/cuteql"
)

func (pg *PgStorage) Categories(ctx context.Context) ([]models.Category, error) {
	categories, tx, err := cuteql.Query[models.Category](ctx, pg.db, `SELECT * FROM categories`)
	if err != nil {
		return nil, err
	}

	return categories, cuteql.Commit(tx)
}
