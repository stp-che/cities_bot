package citiesgame

import (
	_ "embed"
	"math/rand"
	"strings"
)

var (
	//go:embed known_cities.txt
	knownCitiesRaw []byte

	knownCitiesPool *CitiesPool
)

type CitiesPool struct {
	cities []string
}

func NewCitiesPool(cities []string) *CitiesPool {
	return &CitiesPool{cities: cities}
}

func (p *CitiesPool) Includes(city string) bool {
	for _, c := range p.cities {
		if c == city {
			return true
		}
	}

	return false
}

func (p *CitiesPool) GetRandomCity(exceptions []string) (string, bool) {
	n := rand.Intn(len(p.cities)) //nolint: gosec // we do not need strong randomness here
	for i := 0; i < len(p.cities); i++ {
		city := p.cities[(n+i)%len(p.cities)]
		inExceptions := false
		for _, e := range exceptions {
			if e == city {
				inExceptions = true
				break
			}
		}
		if !inExceptions {
			return city, true
		}
	}

	return "", false
}

func KnownCitiesPool() *CitiesPool {
	if knownCitiesPool == nil {
		cities := strings.Split(string(knownCitiesRaw), "\n")
		knownCitiesPool = NewCitiesPool(cities)
	}
	return knownCitiesPool
}
