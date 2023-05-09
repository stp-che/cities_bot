package main

import (
	"context"
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/stp-che/cities_bot/pkg/bot"
	"github.com/stp-che/cities_bot/pkg/log"
)

func main() {
	ctx := context.Background()
	log.Info(ctx, "Cities Bot started")

	bot, err := createBot()
	if err != nil {
		log.Fatal(ctx, fmt.Sprintf("bot init error: %s", err.Error()))
	}

	addHandlers(bot)

	bot.Run(ctx, 0)
}

func createBot() (*bot.Bot, error) {
	cfg := bot.Config{
		Token: os.Getenv("BOT_TOKEN"),
		Debug: os.Getenv("BOT_DEBUG") == "true",
	}

	return bot.New(cfg)
}

func addHandlers(b *bot.Bot) {
	b.AddCommandHandler("play", func(ctx context.Context, msg *tgbotapi.Message) (*tgbotapi.MessageConfig, error) {
		resp := tgbotapi.NewMessage(msg.Chat.ID, "Game started")

		return &resp, nil
	})

	b.SetDefaultHandler(func(ctx context.Context, msg *tgbotapi.Message) (*tgbotapi.MessageConfig, error) {
		resp := tgbotapi.NewMessage(msg.Chat.ID, msg.Text)
		resp.ReplyToMessageID = msg.MessageID

		return &resp, nil
	})
}
