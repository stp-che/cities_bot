package middleware

import (
	"context"
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stp-che/cities_bot/pkg/bot"
	"github.com/stp-che/cities_bot/pkg/log"
	"github.com/stp-che/cities_bot/service/entity/common"
	"github.com/stp-che/cities_bot/service/gateway/telegram"
)

func HandleErrors() func(bot.HandlerFunc) bot.HandlerFunc {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, m *tgbotapi.Message) (*tgbotapi.MessageConfig, error) {
			resp, err := next(ctx, m)
			userErr := telegram.UserError{}
			if errors.As(err, &userErr) {
				errResp := tgbotapi.NewMessage(m.Chat.ID, fmt.Sprintf("Error: %s", userErr.Msg))
				return &errResp, nil
			}

			domainErr := &common.DomainError{}
			if errors.As(err, &domainErr) {
				errResp := tgbotapi.NewMessage(m.Chat.ID, fmt.Sprintf("Error: %s", domainErr.Error()))
				return &errResp, nil
			}

			if err != nil {
				log.Error(ctx, err.Error())
				return nil, nil
			}

			return resp, err
		}
	}
}
