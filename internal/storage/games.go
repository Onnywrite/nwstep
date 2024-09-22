package storage

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/Onnywrite/nwstep/internal/domain/models"
	"github.com/Onnywrite/nwstep/internal/lib/cuteql"
	"github.com/google/uuid"
)

func (pg *PgStorage) IsUserInLobby(ctx context.Context, uid uuid.UUID) bool {
	lobbiesCount, tx, err := cuteql.Get[int](ctx, pg.db, `
	SELECT COUNT(*)
	FROM games_users
	JOIN games ON games_users.game_id = games.game_id
	WHERE user_id = $1 AND games.end_at IS NULL
	`, uid)
	if err != nil {
		return false
	}

	return *lobbiesCount > 0 && cuteql.Commit(tx) == nil
}

func (pg *PgStorage) LinkUserGame(ctx context.Context, ugame models.UserInGame) error {
	tx, err := cuteql.ExecuteSquirreled(ctx, pg.db,
		squirrel.
			Insert("games_users").
			Columns("game_id", "user_id", "health").
			Values(ugame.GameId, ugame.UserId, ugame.Health).
			PlaceholderFormat(squirrel.Dollar),
	)
	if err != nil {
		return err
	}

	return cuteql.Commit(tx)
}

func (pg *PgStorage) CountUsersInGame(ctx context.Context, gameId int) (int, error) {
	count, tx, err := cuteql.Get[int](ctx, pg.db, `
	SELECT COUNT(*) FROM games_users
	WHERE game_id = $1
	`, gameId)
	if err != nil {
		return 0, err
	}

	return *count, cuteql.Commit(tx)
}

func (pg *PgStorage) LobbyGame(ctx context.Context, courseId, needUsers int) (*models.Game, error) {
	game, tx, err := cuteql.Get[models.Game](ctx, pg.db, `
	SELECT * FROM games
	WHERE course_id = $1 AND (
		SELECT COUNT(*) FROM games_users
		WHERE game_id = games.game_id
	) < $2
	`, courseId, needUsers)
	if err != nil {
		return nil, err
	}

	return game, cuteql.Commit(tx)
}

func (pg *PgStorage) SaveGame(ctx context.Context, game models.Game) (*models.Game, error) {
	result, tx, err := cuteql.Get[models.Game](ctx, pg.db, `
	INSERT INTO games (course_id)
	VALUES ($1)
	RETURNING *
	`, game.CourseId)
	if err != nil {
		return nil, err
	}

	return result, cuteql.Commit(tx)
}
