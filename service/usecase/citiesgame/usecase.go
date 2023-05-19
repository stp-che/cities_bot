package citiesgame

import (
	"context"
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
	return "Game started", nil
}

func (u *Usecase) ReceiveMessage(ctx context.Context, message string) (string, error) {
	return message, nil
}
