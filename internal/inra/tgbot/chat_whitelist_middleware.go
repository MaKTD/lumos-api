package tgbot

import (
	"slices"

	tele "gopkg.in/telebot.v4"
)

func ChatWhitelist(chats ...int64) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(ctx tele.Context) error {
			chat := ctx.Chat()
			if chat == nil {
				return nil
			}

			if slices.Contains(chats, chat.ID) {
				return next(ctx)
			}

			return nil
		}
	}
}
