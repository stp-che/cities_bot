package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/stp-che/cities_bot/pkg/log"
)

const (
	DefaultUpdateTimeout = 60
)

type Config struct {
	Token         string
	Debug         bool
	UpdateTimeout int
}

type HandlerFunc func(context.Context, *tgbotapi.Message) (*tgbotapi.MessageConfig, error)

type Bot struct {
	cfg            Config
	botAPI         *tgbotapi.BotAPI
	cmdHandlers    map[string]HandlerFunc
	defaultHandler HandlerFunc
}

func New(cfg Config) (*Bot, error) {
	if cfg.UpdateTimeout == 0 {
		cfg.UpdateTimeout = DefaultUpdateTimeout
	}

	botAPI, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, err
	}

	botAPI.Debug = cfg.Debug

	bot := Bot{
		cfg:         cfg,
		botAPI:      botAPI,
		cmdHandlers: map[string]HandlerFunc{},
	}

	return &bot, nil
}

func (b *Bot) AddCommandHandler(cmd string, h HandlerFunc) {
	b.cmdHandlers[cmd] = h
}

func (b *Bot) SetDefaultHandler(h HandlerFunc) {
	b.defaultHandler = h
}

func (b *Bot) Run(ctx context.Context, offset int) {
	u := tgbotapi.NewUpdate(offset)
	u.Timeout = b.cfg.UpdateTimeout

	updates := b.botAPI.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			msg := update.Message
			msgID := fmt.Sprintf("%d:%d", msg.Chat.ID, msg.MessageID)

			log.Info(ctx, fmt.Sprintf("(%s) [%s] %s", msgID, msg.From.UserName, msg.Text))

			response, err := b.handle(ctx, msg)
			if err != nil {
				log.Warn(ctx, fmt.Sprintf("msg handling error (%s): %s", msgID, err.Error()))
			}

			if response != nil {
				_, err = b.botAPI.Send(response)
				if err != nil {
					log.Warn(ctx, fmt.Sprintf("sending response error (%s): %s", msgID, err.Error()))
				}
			}
		}
	}
}

func (b *Bot) handle(ctx context.Context, msg *tgbotapi.Message) (*tgbotapi.MessageConfig, error) {
	handler := b.defaultHandler

	if msg.IsCommand() {
		if h, ok := b.cmdHandlers[msg.Command()]; ok {
			handler = h
		}
	}

	if handler == nil {
		return nil, nil
	}

	return handler(ctx, msg)
}
