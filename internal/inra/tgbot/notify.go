package tgbot

import (
	"log/slog"
	"strings"

	tele "gopkg.in/telebot.v4"
)

func (r *Bot) ForAdmin(msg string) {
	msg = strings.TrimSpace(msg)
	if msg == "" {
		return
	}

	recipient := tele.ChatID(r.chatID)
	chunks := r.splitTelegramMessage(msg, telegramMaxMessageLen)

	for i, chunk := range chunks {
		chunk = strings.TrimSpace(chunk)
		if chunk == "" {
			continue
		}

		if _, err := r.bot.Send(recipient, chunk); err != nil {
			r.logger.Error(
				"failed to send admin notification message",
				slog.String("err", err.Error()),
				slog.Int64("chat_id", r.chatID),
				slog.Int("chunk_idx", i),
				slog.Int("chunks_total", len(chunks)),
			)
		}
	}
}

func (r *Bot) splitTelegramMessage(s string, maxLen int) []string {
	if maxLen <= 0 {
		return []string{s}
	}
	if len(s) <= maxLen {
		return []string{s}
	}

	rs := []rune(s)
	if len(rs) <= maxLen {
		return []string{s}
	}

	out := make([]string, 0, (len(rs)+maxLen-1)/maxLen)
	for start := 0; start < len(rs); start += maxLen {
		end := start + maxLen
		if end > len(rs) {
			end = len(rs)
		}
		out = append(out, string(rs[start:end]))
	}
	return out
}
