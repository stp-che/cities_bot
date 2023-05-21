package citiesgame

import (
	"context"

	"github.com/google/uuid"
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

func (u *Usecase) Play(ctx context.Context) (*uuid.UUID, string, error) {
	gameUUID := uuid.New()

	return &gameUUID, "Game started", nil
}

func (u *Usecase) ReceiveMessage(ctx context.Context, _ uuid.UUID, message string) (string, error) {
	return message, nil
}

func (u *Usecase) Quit(ctx context.Context, _ uuid.UUID) (string, error) {
	return "Bye!", nil
}
