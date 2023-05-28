package common

import (
	"fmt"

	"github.com/google/uuid"
)

type GameNotFoundError struct {
	engine   string
	gameUUID uuid.UUID
}

func NewGameNotFoundError(engine string, gameUUID uuid.UUID) *GameNotFoundError {
	return &GameNotFoundError{engine: engine, gameUUID: gameUUID}
}

func (e GameNotFoundError) Error() string {
	return fmt.Sprintf("game not found: %s[%s]", e.engine, e.gameUUID)
}
