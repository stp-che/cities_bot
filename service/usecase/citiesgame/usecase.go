package citiesgame

import (
	"context"

	"github.com/google/uuid"
	"github.com/stp-che/cities_bot/service/usecase/session"
)

const (
	gameName = "cities"
)

type Usecase struct{}

func NewUsecase() *Usecase {
	return &Usecase{}
}

func (u *Usecase) Name() string {
	return gameName
}

func (u *Usecase) Play(ctx context.Context) (string, error) {
	s, err := session.GetFromContext(ctx)
	if err != nil {
		return "", err
	}

	if s.Game != nil {
		return "You already have active game", nil
	}

	s.StartGame(u.Name(), uuid.New())

	return "Game started", nil
}

func (u *Usecase) ReceiveMessage(ctx context.Context, message string) (string, error) {
	s, err := session.GetFromContext(ctx)
	if err != nil {
		return "", err
	}

	if s.Game == nil {
		return "You have no current game", nil
	}

	return message, nil
}

func (u *Usecase) Quit(ctx context.Context) (string, error) {
	s, err := session.GetFromContext(ctx)
	if err != nil {
		return "", err
	}

	if s.Game == nil {
		return "You have no current game", nil
	}

	s.QuitGame()

	return "Bye!", nil
}
