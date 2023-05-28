package citiesgame

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGamePlayerTurn(t *testing.T) {
	type testCase struct {
		title       string
		before      func(require.TestingT, *Game)
		city        string
		expectedErr error
		checks      func(require.TestingT, *Game)
	}

	pool := NewCitiesPool([]string{"Foo", "Bar", "Baz"})

	cases := []testCase{
		{
			title:       "when city is not known",
			city:        "Abc",
			expectedErr: ErrCityUnknown,
			checks: func(t require.TestingT, g *Game) {
				require.Empty(t, g.Turns)
				require.False(t, g.IsFinished)
			},
		},
		{
			title: "when city is known",
			city:  "Bar",
			checks: func(t require.TestingT, g *Game) {
				require.Equal(t, 2, len(g.Turns))
				require.Equal(t, "Bar", g.Turns[0])
				require.NotEqual(t, "Bar", g.Turns[1])
				require.False(t, g.IsFinished)
			},
		},
		{
			title: "when city already mentioned",
			city:  "Bar",
			before: func(t require.TestingT, g *Game) {
				require.NoError(t, g.PlayerTurn("Bar"))
			},
			expectedErr: ErrCityMentioned,
			checks: func(t require.TestingT, g *Game) {
				require.False(t, g.IsFinished)
			},
		},
		{
			title: "when all known cities were mentioned",
			city:  "Baz",
			before: func(t require.TestingT, g *Game) {
				g.Turns = []string{"Bar", "Foo"}
			},
			checks: func(t require.TestingT, g *Game) {
				require.Equal(t, []string{"Bar", "Foo", "Baz"}, g.Turns)
				require.True(t, g.IsFinished)
				require.Equal(t, player, g.winner)
			},
		},
		{
			title: "when game is finished",
			city:  "Baz",
			before: func(t require.TestingT, g *Game) {
				g.IsFinished = true
			},
			expectedErr: ErrGameIsFinished,
		},
	}

	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			game := New(pool)

			if tc.before != nil {
				tc.before(t, game)
			}

			errMatcher := func(err error) { require.NoError(t, err) }
			if tc.expectedErr != nil {
				errMatcher = func(err error) { require.ErrorIs(t, err, tc.expectedErr) }
			}

			err := game.PlayerTurn(tc.city)
			errMatcher(err)

			if tc.checks != nil {
				tc.checks(t, game)
			}
		})
	}
}

func TestGamePlayerYields(t *testing.T) {
	type testCase struct {
		title       string
		before      func(require.TestingT, *Game)
		expectedErr error
		checks      func(require.TestingT, *Game)
	}

	pool := NewCitiesPool([]string{"Foo", "Bar", "Baz"})

	cases := []testCase{
		{
			title: "when game is in progress",
			checks: func(t require.TestingT, g *Game) {
				require.Empty(t, g.Turns)
				require.True(t, g.IsFinished)
				require.Equal(t, bot, g.winner)
			},
		},
		{
			title: "when game is finished",
			before: func(t require.TestingT, g *Game) {
				g.IsFinished = true
			},
			expectedErr: ErrGameIsFinished,
		},
	}

	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			game := New(pool)

			if tc.before != nil {
				tc.before(t, game)
			}

			errMatcher := func(err error) { require.NoError(t, err) }
			if tc.expectedErr != nil {
				errMatcher = func(err error) { require.ErrorIs(t, err, tc.expectedErr) }
			}

			err := game.PlayerYields()
			errMatcher(err)

			if tc.checks != nil {
				tc.checks(t, game)
			}
		})
	}
}
