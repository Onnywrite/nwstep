package storage

import (
	"context"

	"github.com/Onnywrite/nwstep/internal/domain/models"
	"github.com/Onnywrite/nwstep/internal/lib/cuteql"

	"github.com/google/uuid"
)

func (pg *PgStorage) Course(ctx context.Context, id int) (*models.Course, error) {
	course, tx, err := cuteql.Get[models.Course](ctx, pg.db, `
	SELECT * FROM courses
	WHERE course_id = $1
	`, id)
	if err != nil {
		return nil, err
	}

	return course, cuteql.Commit(tx)
}

func (pg *PgStorage) Courses(ctx context.Context, categId int64) ([]models.Course, error) {
	courses, tx, err := cuteql.Query[models.Course](ctx, pg.db, `
	SELECT * FROM courses
	WHERE category_id = $1
	ORDER BY min_rating ASC
	`, categId)
	if err != nil {
		return nil, err
	}

	return courses, cuteql.Commit(tx)
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
