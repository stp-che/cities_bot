package session

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/stp-che/cities_bot/pkg/bot"
	"github.com/stp-che/cities_bot/pkg/log"
	esession "github.com/stp-che/cities_bot/service/entity/session"
	"github.com/stp-che/cities_bot/service/usecase/session"
)

type sessionManager interface {
	GetSession(int64) (*esession.Session, error)
	RememberSession(*esession.Session) error
}

func WithSession(sm sessionManager) func(bot.HandlerFunc) bot.HandlerFunc {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, m *tgbotapi.Message) (*tgbotapi.MessageConfig, error) {
			s, err := sm.GetSession(m.Chat.ID)
			if err != nil {
				return nil, err
			}

			if s == nil {
				s = esession.New(m.Chat.ID)
			}

			defer func() {
				err := sm.RememberSession(s)
				if err != nil {
					log.Warn(ctx, fmt.Sprintf("saving session error: %s", err.Error()))
				}
			}()

			ctx = session.NewContext(ctx, s)

			return next(ctx, m)
		}
	}
}
