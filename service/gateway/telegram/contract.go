package telegram

import (
	"context"

	"github.com/google/uuid"
)

//go:generate mockgen -source=contract.go -destination=contract_mock.go -package=telegram

type GameEngine interface {
	Name() string
	Play(context.Context) (*uuid.UUID, string, error)
	ReceiveMessage(context.Context, uuid.UUID, string) (string, bool, error)
	Yield(context.Context, uuid.UUID) (string, bool, error)
	Quit(context.Context, uuid.UUID) (string, error)
}
