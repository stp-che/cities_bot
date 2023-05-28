package citiesgame

import (
	"context"

	"github.com/google/uuid"
	"github.com/stp-che/cities_bot/service/entity/citiesgame"
)

type MemRepo struct {
	games map[uuid.UUID]*citiesgame.Game
}

func NewMemRepo() *MemRepo {
	return &MemRepo{
		games: map[uuid.UUID]*citiesgame.Game{},
	}
}

func (r *MemRepo) Get(_ context.Context, gameUUID uuid.UUID) (*citiesgame.Game, error) {
	return r.games[gameUUID], nil
}

func (r *MemRepo) Save(_ context.Context, game *citiesgame.Game) error {
	if game.UUID == uuid.Nil {
		game.UUID = uuid.New()
	}

	r.games[game.UUID] = game

	return nil
}

func (r *MemRepo) Delete(_ context.Context, gameUUID uuid.UUID) error {
	delete(r.games, gameUUID)
	return nil
}
