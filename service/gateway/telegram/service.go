package telegram

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	esession "github.com/stp-che/cities_bot/service/entity/session"
	"github.com/stp-che/cities_bot/service/usecase/session"
)

type Service struct {
	games map[string]GameEngine
}

func NewService(games []GameEngine) *Service {
	gamesMap := map[string]GameEngine{}
	for _, g := range games {
		gamesMap[g.Name()] = g
	}
	return &Service{games: gamesMap}
}

func (s *Service) Play(ctx context.Context, msg *tgbotapi.Message) (*tgbotapi.MessageConfig, error) {
	engine, err := s.selectGame(msg.CommandArguments())
	if err != nil {
		return nil, err
	}

	sess, err := session.GetFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if sess.Game != nil {
		return nil, ErrGameAlreadyStarted
	}

	gameUUID, res, err := engine.Play(ctx)
	if err != nil {
		return nil, err
	}

	sess.StartGame(engine.Name(), *gameUUID)

	resp := tgbotapi.NewMessage(msg.Chat.ID, res)

	return &resp, nil
}

func (s *Service) Default(ctx context.Context, msg *tgbotapi.Message) (*tgbotapi.MessageConfig, error) {
	return s.withCurrentGame(ctx,
		func(sess *esession.Session, engine GameEngine, gameUUID uuid.UUID) (*tgbotapi.MessageConfig, error) {
			res, _, err := engine.ReceiveMessage(ctx, sess.Game.UUID, msg.Text)
			if err != nil {
				return nil, err
			}

			resp := tgbotapi.NewMessage(msg.Chat.ID, res)

			return &resp, nil
		},
	)
}

func (s *Service) Quit(ctx context.Context, msg *tgbotapi.Message) (*tgbotapi.MessageConfig, error) {
	return s.withCurrentGame(ctx,
		func(sess *esession.Session, engine GameEngine, gameUUID uuid.UUID) (*tgbotapi.MessageConfig, error) {
			res, err := engine.Quit(ctx, sess.Game.UUID)
			if err != nil {
				return nil, err
			}

			sess.QuitGame()

			resp := tgbotapi.NewMessage(msg.Chat.ID, res)

			return &resp, nil
		},
	)
}

func (s *Service) withCurrentGame(
	ctx context.Context, handle func(*esession.Session, GameEngine, uuid.UUID) (*tgbotapi.MessageConfig, error),
) (*tgbotapi.MessageConfig, error) {
	sess, err := session.GetFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if sess.Game == nil {
		return nil, ErrGameNotStarted
	}

	game, exists := s.games[sess.Game.Name]
	if !exists {
		sess.QuitGame() // clear bad game from session
		return nil, ErrGameNotStarted
	}

	return handle(sess, game, sess.Game.UUID)
}

func (s *Service) selectGame(name string) (GameEngine, error) {
	if name == "" {
		return nil, ErrGameNameNotSpecified
	}

	game, exists := s.games[name]
	if !exists {
		return nil, ErrGameDoesNotExist
	}

	return game, nil
}
