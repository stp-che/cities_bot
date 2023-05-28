package citiesgame

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stp-che/cities_bot/service/entity/citiesgame"
	"github.com/stp-che/cities_bot/service/entity/citiesgame/mocks"
	"github.com/stp-che/cities_bot/service/entity/common"
	"github.com/stretchr/testify/require"
)

func TestUsecasePlay(t *testing.T) {
	type testCase struct {
		title  string
		before func(*testCase)
		check  func(require.TestingT, error)

		gameRepo *mocks.MockRepository
	}

	whateverError := errors.New("whatever")

	cases := []*testCase{
		{
			title: "when no error happened",
			before: func(tc *testCase) {
				tc.gameRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
			},
			check: func(t require.TestingT, err error) {
				require.NoError(t, err)
			},
		},
		{
			title: "when game repo returns error",
			before: func(tc *testCase) {
				tc.gameRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(whateverError)
			},
			check: func(t require.TestingT, err error) {
				require.ErrorIs(t, err, whateverError)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			tc.gameRepo = mocks.NewMockRepository(ctrl)

			if tc.before != nil {
				tc.before(tc)
			}

			ctx := context.Background()

			u := NewUsecase(WithGameRepo(tc.gameRepo))
			_, _, err := u.Play(ctx)

			if tc.check != nil {
				tc.check(t, err)
			}
		})
	}
}

func TestUsecaseReceiveMessage(t *testing.T) {
	type testCase struct {
		title  string
		msg    string
		before func(*testCase)
		check  func(require.TestingT, string, bool, error)

		gameRepo *mocks.MockRepository
	}

	whateverError := errors.New("whatever")
	gameUUID := uuid.New()

	cases := []*testCase{
		{
			title: "when no error happened and game continues",
			msg:   "London",
			before: func(tc *testCase) {
				tc.gameRepo.EXPECT().Get(gomock.Any(), gameUUID).Return(testGame([]string{"London", "Berlin"}), nil)
				tc.gameRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
			},
			check: func(t require.TestingT, res string, isFinished bool, err error) {
				require.NoError(t, err)
				require.Equal(t, "Berlin", res)
				require.False(t, isFinished)
			},
		},
		{
			title: "when no error happened and game finishes",
			msg:   "London",
			before: func(tc *testCase) {
				tc.gameRepo.EXPECT().Get(gomock.Any(), gameUUID).Return(testGame([]string{"London"}), nil)
				tc.gameRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
			},
			check: func(t require.TestingT, res string, isFinished bool, err error) {
				require.NoError(t, err)
				require.NotZero(t, res)
				require.True(t, isFinished)
			},
		},
		{
			title: "when game repo returns error",
			msg:   "London",
			before: func(tc *testCase) {
				tc.gameRepo.EXPECT().Get(gomock.Any(), gameUUID).Return(testGame([]string{"London", "Berlin"}), nil)
				tc.gameRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(whateverError)
			},
			check: func(t require.TestingT, _ string, _ bool, err error) {
				require.ErrorIs(t, err, whateverError)
			},
		},
		{
			title: "when game does not exist",
			msg:   "London",
			before: func(tc *testCase) {
				tc.gameRepo.EXPECT().Get(gomock.Any(), gameUUID).Return(nil, nil)
			},
			check: func(t require.TestingT, _ string, _ bool, err error) {
				notFoundErr := &common.GameNotFoundError{}
				require.ErrorAs(t, err, &notFoundErr)
				require.Equal(t, common.NewGameNotFoundError(gameName, gameUUID), notFoundErr)
			},
		},
		{
			title: "when game returns error",
			msg:   "foo", // we pass unknown city in order to get error from the Game
			before: func(tc *testCase) {
				tc.gameRepo.EXPECT().Get(gomock.Any(), gameUUID).Return(testGame([]string{"London", "Berlin"}), nil)
			},
			check: func(t require.TestingT, _ string, _ bool, err error) {
				domainErr := &common.DomainError{}
				require.ErrorAs(t, err, &domainErr)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			tc.gameRepo = mocks.NewMockRepository(ctrl)

			if tc.before != nil {
				tc.before(tc)
			}

			ctx := context.Background()

			u := NewUsecase(WithGameRepo(tc.gameRepo))
			res, isFinished, err := u.ReceiveMessage(ctx, gameUUID, tc.msg)

			if tc.check != nil {
				tc.check(t, res, isFinished, err)
			}
		})
	}
}

func testGame(cities []string) *citiesgame.Game {
	citiesPool := citiesgame.NewCitiesPool(cities)
	return citiesgame.New(citiesPool)
}
