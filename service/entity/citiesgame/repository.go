package citiesgame

import (
	"context"

	"github.com/google/uuid"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go -package=mocks

type Repository interface {
	Save(context.Context, *Game) error
	Get(context.Context, uuid.UUID) (*Game, error)
	Delete(context.Context, uuid.UUID) error
}
