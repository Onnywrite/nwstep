package storage

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/Onnywrite/nwstep/internal/domain/models"
	"github.com/Onnywrite/nwstep/internal/lib/cuteql"
)

func (pg *PgStorage) RandomQuestions(ctx context.Context, courseId, count int) ([]models.Question, error) {
	questions, tx, err := cuteql.Query[models.Question](ctx, pg.db, `
	SELECT * FROM questions
	WHERE course_id = $1
	ORDER BY RANDOM()
	LIMIT $2
	`, courseId, count)
	if err != nil {
		return nil, err
	}

	return questions, cuteql.Commit(tx)
}

func (pg *PgStorage) SaveGameQuestions(ctx context.Context, questions ...models.GameQuestion) error {
	insertBuilder := squirrel.Insert("games_questions").
		Columns("game_id", "question_id", "number")

	for _, question := range questions {
		insertBuilder = insertBuilder.Values(question.GameId, question.QuestionId, question.Number)
	}

	tx, err := cuteql.ExecuteSquirreled(ctx, pg.db, insertBuilder.PlaceholderFormat(squirrel.Dollar))
	if err != nil {
		return err
	}

	return cuteql.Commit(tx)
}

func (pg *PgStorage) Answers(ctx context.Context, questionId int) ([]models.Answer, error) {
	answers, tx, err := cuteql.Query[models.Answer](ctx, pg.db, `
	SELECT * FROM answers
	WHERE question_id = $1
	`, questionId)
	if err != nil {
		return nil, err
	}

	return answers, cuteql.Commit(tx)
}

func (pg *PgStorage) Game(ctx context.Context, gameId int) (*models.Game, error) {
	game, err := pg.getGameWhere(ctx, squirrel.Eq{"game_id": gameId})
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (pg *PgStorage) getGameWhere(ctx context.Context, where squirrel.Sqlizer,
) (*models.Game, error) {
	game, tx, err := cuteql.GetSquirreled[models.Game](ctx, pg.db,
		squirrel.
			Select("*").
			From("games").
			Where(where).PlaceholderFormat(squirrel.Dollar),
	)
	if err != nil {
		return nil, err
	}

	return game, cuteql.Commit(tx)
}

func (pg *PgStorage) UpdateGame(ctx context.Context, gameId int, newValues map[string]any) error {
	if len(newValues) == 0 {
		return fmt.Errorf("no fields to update: %w", cuteql.ErrEmptyResult)
	}

	if _, ok := newValues["game_id"]; ok {
		return fmt.Errorf("game_id must not be changed: %w", cuteql.ErrInternal)
	}

	tx, err := cuteql.ExecuteSquirreled(ctx, pg.db,
		squirrel.
			Update("games").
			SetMap(newValues).
			Where("game_id = ?", gameId).PlaceholderFormat(squirrel.Dollar),
	)
	if err != nil {
		return err
	}

	return cuteql.Commit(tx)
}

func (pg *PgStorage) QuestionByNumber(ctx context.Context, gameId, questionNumber int) (*models.Question, error) {
	question, tx, err := cuteql.Get[models.Question](ctx, pg.db, `
	SELECT questions.* FROM games_questions
	JOIN questions ON questions.question_id = games_questions.question_id
	WHERE games_questions.game_id = $1 AND games_questions.number = $2
	`, gameId, questionNumber)
	if err != nil {
		return nil, err
	}

	return question, cuteql.Commit(tx)
}
