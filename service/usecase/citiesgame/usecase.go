package citiesgame

import (
	"context"

	"github.com/google/uuid"
	"github.com/stp-che/cities_bot/service/entity/citiesgame"
	"github.com/stp-che/cities_bot/service/entity/common"
)

const (
	gameName = "cities"
)

type Usecase struct {
	gameRepo citiesgame.Repository
}

func NewUsecase(opts ...Option) *Usecase {
	u := &Usecase{}

	for _, o := range opts {
		o(u)
	}

	return u
}

func (u *Usecase) Name() string {
	return gameName
}

func (u *Usecase) Play(ctx context.Context) (*uuid.UUID, string, error) {
	game := citiesgame.New(citiesgame.KnownCitiesPool())

	err := u.gameRepo.Save(ctx, game)
	if err != nil {
		return nil, "", err
	}

	return &game.UUID, game.Greeting(), nil
}

func (u *Usecase) ReceiveMessage(ctx context.Context, gameUUID uuid.UUID, message string) (string, bool, error) {
	game, err := u.gameRepo.Get(ctx, gameUUID)
	if err != nil {
		return "", false, err
	}

	if game == nil {
		return "", false, common.NewGameNotFoundError(gameName, gameUUID)
	}

	err = game.PlayerTurn(message)
	if err != nil {
		return "", false, err
	}

	err = u.gameRepo.Save(ctx, game)
	if err != nil {
		return "", false, err
	}

	resp := game.LastTurn()
	if game.IsFinished {
		resp = game.Result()
	}

	return resp, game.IsFinished, nil
}

func (u *Usecase) Yield(ctx context.Context, gameUUID uuid.UUID) (string, bool, error) {
	game, err := u.gameRepo.Get(ctx, gameUUID)
	if err != nil {
		return "", false, err
	}

	if game == nil {
		return "", false, common.NewGameNotFoundError(gameName, gameUUID)
	}

	err = game.PlayerYields()
	if err != nil {
		return "", false, err
	}

	err = u.gameRepo.Save(ctx, game)
	if err != nil {
		return "", false, err
	}

	return game.Result(), game.IsFinished, nil
}

func (u *Usecase) Quit(ctx context.Context, gameUUID uuid.UUID) (string, error) {
	err := u.gameRepo.Delete(ctx, gameUUID)
	if err != nil {
		return "", err
	}

	return "Bye!", nil
}
