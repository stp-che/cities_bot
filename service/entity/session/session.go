package session

import (
	"github.com/google/uuid"
)

type Session struct {
	ChatID int64
	Game   *Game
}

type Game struct {
	Name       string
	UUID       uuid.UUID
	IsFinished bool
}

func New(chatID int64) *Session {
	return &Session{
		ChatID: chatID,
	}
}

func (s *Session) StartGame(name string, uuid uuid.UUID) {
	s.Game = &Game{Name: name, UUID: uuid}
}

func (s *Session) QuitGame() {
	s.Game = nil
}
