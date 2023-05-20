package session

import (
	"context"
	"errors"

	"github.com/stp-che/cities_bot/service/entity/session"
)

type ctxSessionKeyType struct{}

var (
	ctxSessionKey     = ctxSessionKeyType{}
	errInvalidSession = errors.New("invalid session")
)

func NewContext(ctx context.Context, s *session.Session) context.Context {
	return context.WithValue(ctx, ctxSessionKey, s)
}

func GetFromContext(ctx context.Context) (*session.Session, error) {
	v := ctx.Value(ctxSessionKey)

	if v == nil {
		return nil, nil
	}

	s, ok := v.(*session.Session)
	if !ok {
		return nil, errInvalidSession
	}

	return s, nil
}
