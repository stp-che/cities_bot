package citiesgame

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCitiesPoolIncludes(t *testing.T) {
	pool := NewCitiesPool([]string{"Foo", "Bar"})

	require.True(t, pool.Includes("Foo"))
	require.True(t, pool.Includes("Bar"))
	require.False(t, pool.Includes("Baz"))
}

func TestCitiesPoolGetRandomCity(t *testing.T) {
	pool := NewCitiesPool([]string{"Foo", "Bar"})

	city, ok := pool.GetRandomCity([]string{})
	require.True(t, ok)
	require.True(t, city == "Foo" || city == "Bar")

	city, ok = pool.GetRandomCity([]string{"Foo"})
	require.True(t, ok)
	require.Equal(t, "Bar", city)

	city, ok = pool.GetRandomCity([]string{"Foo", "Bar"})
	require.False(t, ok)
	require.Zero(t, city)
}
