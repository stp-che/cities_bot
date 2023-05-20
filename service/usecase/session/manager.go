package session

import (
	"github.com/stp-che/cities_bot/service/entity/session"
)

type Manager struct {
	sessionRepo session.Repository
}

func NewManager(repo session.Repository) *Manager {
	return &Manager{sessionRepo: repo}
}

func (m *Manager) GetSession(chatID int64) (*session.Session, error) {
	s, err := m.sessionRepo.GetByChatID(chatID)
	if err != nil {
		return nil, err
	}

	if s == nil {
		s = session.New(chatID)
	}

	return s, nil
}

func (m *Manager) RememberSession(s *session.Session) error {
	return m.sessionRepo.Save(s)
}
