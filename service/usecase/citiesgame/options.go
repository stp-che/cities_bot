package citiesgame

import "github.com/stp-che/cities_bot/service/entity/citiesgame"

type Option func(*Usecase)

func WithGameRepo(repo citiesgame.Repository) Option {
	return func(u *Usecase) {
		u.gameRepo = repo
	}
}
