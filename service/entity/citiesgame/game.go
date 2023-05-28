package citiesgame

import (
	"github.com/google/uuid"
	"github.com/stp-che/cities_bot/service/entity/common"
)

type Participant string

const (
	bot    Participant = "bot"
	player Participant = "player"
)

var (
	ErrCityUnknown    = common.NewDomainError("city unknown")
	ErrCityMentioned  = common.NewDomainError("city already mentioned")
	ErrGameIsFinished = common.NewDomainError("game is finished")
)

const (
	greeting      = "Game started"
	resPlayerWins = "Bot yields\nGame ended! You won!"
	resBotWins    = "You yield\nGame ended"
)

type Game struct {
	UUID       uuid.UUID
	Turns      []string
	IsFinished bool
	CitiesPool *CitiesPool
	winner     Participant
}

func New(citiesPool *CitiesPool) *Game {
	return &Game{
		UUID:       uuid.New(),
		CitiesPool: citiesPool,
	}
}

func (g *Game) PlayerTurn(city string) error {
	if g.IsFinished {
		return ErrGameIsFinished
	}

	for _, c := range g.Turns {
		if c == city {
			return ErrCityMentioned
		}
	}

	if !g.CitiesPool.Includes(city) {
		return ErrCityUnknown
	}

	g.Turns = append(g.Turns, city)

	botChoice, ok := g.CitiesPool.GetRandomCity(g.Turns)
	if !ok {
		g.finish(player)
		return nil
	}

	g.Turns = append(g.Turns, botChoice)

	return nil
}

func (g *Game) PlayerYields() error {
	if g.IsFinished {
		return ErrGameIsFinished
	}

	g.finish(bot)

	return nil
}

func (g *Game) Greeting() string {
	return greeting
}

func (g *Game) Result() string {
	if !g.IsFinished {
		return ""
	}

	switch g.winner {
	case bot:
		return resBotWins
	case player:
		return resPlayerWins
	}

	return ""
}

func (g *Game) LastTurn() string {
	if len(g.Turns) == 0 {
		return ""
	}

	return g.Turns[len(g.Turns)-1]
}

func (g *Game) finish(winner Participant) {
	g.IsFinished = true
	g.winner = winner
}
