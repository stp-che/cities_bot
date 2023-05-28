package main

import (
	"context"
	"fmt"
	"os"

	"github.com/stp-che/cities_bot/pkg/bot"
	"github.com/stp-che/cities_bot/pkg/log"
	"github.com/stp-che/cities_bot/service/gateway/telegram"
	tgmd "github.com/stp-che/cities_bot/service/gateway/telegram/middleware"
	citiesgamerepo "github.com/stp-che/cities_bot/service/repository/citiesgame"
	sessrepo "github.com/stp-che/cities_bot/service/repository/session"
	"github.com/stp-che/cities_bot/service/usecase/citiesgame"
	"github.com/stp-che/cities_bot/service/usecase/session"
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
	cfg     config
	bot     *bot.Bot
	sessMgr *session.Manager

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

	a.TgHandler = telegram.NewService(
		[]telegram.GameEngine{
			citiesgame.NewUsecase(citiesgame.WithGameRepo(citiesgamerepo.NewMemRepo())),
		},
	)

	a.addBotHandlers()

	return nil
}

func (a *App) Run(ctx context.Context) {
	a.bot.Run(ctx, 0)
}

func (a *App) addBotHandlers() {
	withSession := tgmd.WithSession(a.sessionManager())
	handleErrors := tgmd.HandleErrors()

	m := func(h bot.HandlerFunc) bot.HandlerFunc {
		return handleErrors(withSession(h))
	}

	a.bot.AddCommandHandler("play", m(a.TgHandler.Play))
	a.bot.AddCommandHandler("yield", m(a.TgHandler.Yield))
	a.bot.AddCommandHandler("quit", m(a.TgHandler.Quit))

	a.bot.SetDefaultHandler(m(a.TgHandler.Default))
}

func (a *App) sessionManager() *session.Manager {
	if a.sessMgr == nil {
		repo := sessrepo.New()
		a.sessMgr = session.NewManager(repo)
	}

	return a.sessMgr
}
