package telegram

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Game interface {
	Name() string
	Play(context.Context) (string, error)
	ReceiveMessage(context.Context, string) (string, error)
	Quit(context.Context) (string, error)
}

type Service struct {
	games map[string]Game
}

func NewService(games []Game) *Service {
	gamesMap := map[string]Game{}
	for _, g := range games {
		gamesMap[g.Name()] = g
	}
	return &Service{games: gamesMap}
}

func (s *Service) Play(ctx context.Context, msg *tgbotapi.Message) (*tgbotapi.MessageConfig, error) {
	res, err := s.getCurrentGame(ctx).Play(ctx)
	if err != nil {
		return nil, err
	}

	resp := tgbotapi.NewMessage(msg.Chat.ID, res)

	return &resp, nil
}

func (s *Service) Default(ctx context.Context, msg *tgbotapi.Message) (*tgbotapi.MessageConfig, error) {
	res, err := s.getCurrentGame(ctx).ReceiveMessage(ctx, msg.Text)
	if err != nil {
		return nil, err
	}

	resp := tgbotapi.NewMessage(msg.Chat.ID, res)

	return &resp, nil
}

func (s Service) Quit(ctx context.Context, msg *tgbotapi.Message) (*tgbotapi.MessageConfig, error) {
	res, err := s.getCurrentGame(ctx).Quit(ctx)
	if err != nil {
		return nil, err
	}

	resp := tgbotapi.NewMessage(msg.Chat.ID, res)

	return &resp, nil
}

func (s *Service) getCurrentGame(ctx context.Context) Game {
	return s.games["cities"]
}
