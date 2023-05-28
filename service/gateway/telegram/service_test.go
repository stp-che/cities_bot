package telegram

import (
	"context"
	"errors"
	"strings"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stp-che/cities_bot/pkg/bot"
	esession "github.com/stp-che/cities_bot/service/entity/session"
	"github.com/stp-che/cities_bot/service/usecase/session"
	"github.com/stretchr/testify/require"
)

func TestServicePlay(t *testing.T) {
	errWhatever := errors.New("whatever")
	gameUUID := uuid.New()

	cases := []*testCase{
		{
			title: "when game name not specified",
			msg:   cmdMsg("/play"),
			checks: func(t require.TestingT, _ *testCase, _ *tgbotapi.MessageConfig, err error) {
				require.ErrorIs(t, err, ErrGameNameNotSpecified)
			},
		},
		{
			title: "when there is no current game in session",
			msg:   cmdMsg("/play", "foo"),
			before: func(tc *testCase) {
				tc.engines["foo"].EXPECT().Play(gomock.Any()).Return(&gameUUID, "answer", nil)
			},
			checks: func(t require.TestingT, tc *testCase, mc *tgbotapi.MessageConfig, err error) {
				require.NoError(t, err)
				require.Equal(t, "answer", mc.Text)
				require.Equal(t, &esession.Game{Name: "foo", UUID: gameUUID}, tc.session.Game)
			},
		},
		{
			title: "when there is current unfinished game in session",
			msg:   cmdMsg("/play", "foo"),
			before: func(tc *testCase) {
				tc.session.StartGame("foo", uuid.New())
			},
			checks: func(t require.TestingT, _ *testCase, mc *tgbotapi.MessageConfig, err error) {
				require.ErrorIs(t, err, ErrGameAlreadyStarted)
			},
		},
		{
			title: "when there is finished game in session",
			msg:   cmdMsg("/play", "foo"),
			before: func(tc *testCase) {
				tc.session.StartGame("foo", uuid.New())
				tc.session.Game.IsFinished = true
				tc.engines["foo"].EXPECT().Play(gomock.Any()).Return(&gameUUID, "answer", nil)
			},
			checks: func(t require.TestingT, tc *testCase, mc *tgbotapi.MessageConfig, err error) {
				require.NoError(t, err)
				require.Equal(t, "answer", mc.Text)
				require.Equal(t, &esession.Game{Name: "foo", UUID: gameUUID}, tc.session.Game)
			},
		},
		{
			title: "when there is no game with given name",
			msg:   cmdMsg("/play", "baz"),
			checks: func(t require.TestingT, _ *testCase, _ *tgbotapi.MessageConfig, err error) {
				require.ErrorIs(t, err, ErrGameDoesNotExist)
			},
		},
		{
			title: "when usecase returns an error",
			msg:   cmdMsg("/play", "foo"),
			before: func(tc *testCase) {
				tc.engines["foo"].EXPECT().Play(gomock.Any()).Return(nil, "", errWhatever)
			},
			checks: func(t require.TestingT, _ *testCase, mc *tgbotapi.MessageConfig, err error) {
				require.ErrorIs(t, err, errWhatever)
				require.Nil(t, mc)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			tc.run(t, func(s *Service) bot.HandlerFunc {
				return s.Play
			})
		})
	}
}

func TestServiceQuit(t *testing.T) {
	gameUUID := uuid.New()

	cases := []*testCase{
		{
			title: "when there is no game in session",
			msg:   cmdMsg("/quit"),
			checks: func(tt require.TestingT, _ *testCase, _ *tgbotapi.MessageConfig, err error) {
				require.ErrorIs(t, err, ErrGameNotStarted)
			},
		},
		{
			title: "when game started",
			msg:   cmdMsg("/quit"),
			before: func(tc *testCase) {
				tc.session.StartGame("foo", gameUUID)
				tc.engines["foo"].EXPECT().Quit(gomock.Any(), gameUUID).Return("answer", nil)
			},
			checks: func(tt require.TestingT, tc *testCase, mc *tgbotapi.MessageConfig, err error) {
				require.NoError(t, err)
				require.Equal(t, "answer", mc.Text)
				require.Nil(t, tc.session.Game, "game should be removed from session")
			},
		},
		{
			title: "when there is a game with unknown engine in session",
			msg:   cmdMsg("/quit"),
			before: func(tc *testCase) {
				tc.session.StartGame("bar", gameUUID)
			},
			checks: func(tt require.TestingT, _ *testCase, _ *tgbotapi.MessageConfig, err error) {
				require.ErrorIs(t, err, ErrGameNotStarted)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			tc.run(t, func(s *Service) bot.HandlerFunc {
				return s.Quit
			})
		})
	}
}

func TestServiceYield(t *testing.T) {
	gameUUID := uuid.New()

	cases := []*testCase{
		{
			title: "when there is no game in session",
			msg:   cmdMsg("/yield"),
			checks: func(tt require.TestingT, _ *testCase, _ *tgbotapi.MessageConfig, err error) {
				require.ErrorIs(t, err, ErrGameNotStarted)
			},
		},
		{
			title: "when game started",
			msg:   cmdMsg("/yield"),
			before: func(tc *testCase) {
				tc.session.StartGame("foo", gameUUID)
				tc.engines["foo"].EXPECT().Yield(gomock.Any(), gameUUID).Return("answer", true, nil)
			},
			checks: func(tt require.TestingT, tc *testCase, mc *tgbotapi.MessageConfig, err error) {
				require.NoError(t, err)
				require.Equal(t, "answer", mc.Text)
				require.True(t, tc.session.Game.IsFinished)
			},
		},
		{
			title: "when there is a game with unknown engine in session",
			msg:   cmdMsg("/yield"),
			before: func(tc *testCase) {
				tc.session.StartGame("bar", gameUUID)
			},
			checks: func(tt require.TestingT, _ *testCase, _ *tgbotapi.MessageConfig, err error) {
				require.ErrorIs(t, err, ErrGameNotStarted)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			tc.run(t, func(s *Service) bot.HandlerFunc {
				return s.Yield
			})
		})
	}
}

func TestServiceDefault(t *testing.T) {
	gameUUID := uuid.New()

	cases := []*testCase{
		{
			title: "when there is no game in session",
			msg:   msg("something"),
			checks: func(tt require.TestingT, _ *testCase, _ *tgbotapi.MessageConfig, err error) {
				require.ErrorIs(t, err, ErrGameNotStarted)
			},
		},
		{
			title: "when game continues",
			msg:   msg("something"),
			before: func(tc *testCase) {
				tc.session.StartGame("foo", gameUUID)
				tc.engines["foo"].EXPECT().ReceiveMessage(gomock.Any(), gameUUID, "something").Return("answer", false, nil)
			},
			checks: func(tt require.TestingT, tc *testCase, mc *tgbotapi.MessageConfig, err error) {
				require.NoError(t, err)
				require.Equal(t, "answer", mc.Text)
				require.False(t, tc.session.Game.IsFinished)
			},
		},
		{
			title: "when game finishes",
			msg:   msg("something"),
			before: func(tc *testCase) {
				tc.session.StartGame("foo", gameUUID)
				tc.engines["foo"].EXPECT().ReceiveMessage(gomock.Any(), gameUUID, "something").Return("answer", true, nil)
			},
			checks: func(tt require.TestingT, tc *testCase, mc *tgbotapi.MessageConfig, err error) {
				require.NoError(t, err)
				require.Equal(t, "answer", mc.Text)
				require.True(t, tc.session.Game.IsFinished)
			},
		},
		{
			title: "when there is a game with unknown engine in session",
			msg:   msg("something"),
			before: func(tc *testCase) {
				tc.session.StartGame("bar", gameUUID)
			},
			checks: func(tt require.TestingT, _ *testCase, _ *tgbotapi.MessageConfig, err error) {
				require.ErrorIs(t, err, ErrGameNotStarted)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			tc.run(t, func(s *Service) bot.HandlerFunc {
				return s.Default
			})
		})
	}
}

type testCase struct {
	title  string
	before func(*testCase)
	msg    *tgbotapi.Message
	checks func(require.TestingT, *testCase, *tgbotapi.MessageConfig, error)

	engines map[string]*MockGameEngine
	session *esession.Session
}

func (tc *testCase) run(t *testing.T, handler func(*Service) bot.HandlerFunc) {
	ctrl := gomock.NewController(t)
	tc.engines = map[string]*MockGameEngine{
		"foo": testGameEngine(ctrl, "foo"),
	}
	tc.session = esession.New(1)

	service := NewService([]GameEngine{tc.engines["foo"]})

	if tc.before != nil {
		tc.before(tc)
	}

	ctx := session.NewContext(context.Background(), tc.session)

	res, err := handler(service)(ctx, tc.msg)

	if tc.checks != nil {
		tc.checks(t, tc, res, err)
	}
}

func testGameEngine(ctrl *gomock.Controller, name string) *MockGameEngine {
	engine := NewMockGameEngine(ctrl)
	engine.EXPECT().Name().Return(name).AnyTimes()

	return engine
}

func cmdMsg(cmd string, cmdArgs ...string) *tgbotapi.Message {
	return &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: 1},
		Text: strings.Join(append([]string{cmd}, cmdArgs...), " "),
		Entities: []tgbotapi.MessageEntity{
			{
				Type:   "bot_command",
				Offset: 0,
				Length: len(cmd),
			},
		},
	}
}

func msg(text string) *tgbotapi.Message {
	return &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: 1},
		Text: text,
	}
}
