package main

import (
	"context"
	"fmt"
	"os"

	"github.com/stp-che/cities_bot/pkg/bot"
	"github.com/stp-che/cities_bot/pkg/log"
	"github.com/stp-che/cities_bot/service/gateway/telegram"
	"github.com/stp-che/cities_bot/service/usecase/citiesgame"
)

func main() {
	ctx := context.Background()
	log.Info(ctx, "Cities Bot started")

	app := NewApp()

	err := app.Init()
	if err != nil {
		log.Fatal(ctx, fmt.Sprintf("app init error: %s", err.Error()))
	}

	app.Run(ctx)
}

type config struct {
	bot bot.Config
}

type App struct {
	cfg config
	bot *bot.Bot

	TgHandler *telegram.Service
}

func NewApp() *App {
	return &App{
		cfg: config{
			bot: bot.Config{
				Token: os.Getenv("BOT_TOKEN"),
				Debug: os.Getenv("BOT_DEBUG") == "true",
			},
		},
	}
}

func (a *App) Init() error {
	var err error

	a.bot, err = bot.New(a.cfg.bot)
	if err != nil {
		return fmt.Errorf("bot init error: %w", err)
	}

	a.TgHandler = telegram.NewService([]telegram.Game{citiesgame.NewUsecase()})

	a.addBotHandlers()

	return nil
}

func (a *App) Run(ctx context.Context) {
	a.bot.Run(ctx, 0)
}

func (a *App) addBotHandlers() {
	a.bot.AddCommandHandler("play", a.TgHandler.Play)

	a.bot.SetDefaultHandler(a.TgHandler.Default)
}
