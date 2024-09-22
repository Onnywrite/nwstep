package storage

import (
	"context"

	"github.com/Onnywrite/nwstep/internal/domain/models"
	"github.com/Onnywrite/nwstep/internal/lib/cuteql"
)

func (pg *PgStorage) DeleteCourse(ctx context.Context, courseId int) error {
	tx, err := cuteql.Execute(ctx, pg.db, `
	DELETE FROM courses
	WHERE course_id = $1
	`, courseId)
	if err != nil {
		return err
	}

	return cuteql.Commit(tx)
}

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

func (pg *PgStorage) Courses(ctx context.Context, categId int) ([]models.Course, error) {
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
