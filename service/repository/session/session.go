package session

import (
	"github.com/stp-che/cities_bot/service/entity/session"
)

type Repository struct {
	data map[int64]*session.Session
}

func New() *Repository {
	return &Repository{
		data: make(map[int64]*session.Session),
	}
}

func (r *Repository) GetByChatID(chatID int64) (*session.Session, error) {
	return r.data[chatID], nil
}

func (r *Repository) Save(s *session.Session) error {
	r.data[s.ChatID] = s
	return nil
}
