// Package telegram provides Telegram Bot API integration for sending notifications.
package telegram

import (
	"context"
	"fmt"
	"log/slog"

	"gabe565.com/domain-watch/internal/config"
	"gabe565.com/domain-watch/internal/util"
	"github.com/go-telegram/bot"
)

// Telegram implements the Integration interface for Telegram notifications.
type Telegram struct {
	ChatID int64
	Bot    *bot.Bot
}

// Name returns the integration name.
func (t *Telegram) Name() string { return "Telegram" }

// Setup initializes the Telegram integration with the provided configuration.
func (t *Telegram) Setup(ctx context.Context, conf *config.Config) error {
	if t.ChatID = conf.TelegramChat; t.ChatID == 0 {
		return fmt.Errorf("telegram %w: chat ID", util.ErrNotConfigured)
	}

	return t.Login(ctx, conf.TelegramToken)
}

// Login authenticates with the Telegram Bot API using the provided token.
func (t *Telegram) Login(ctx context.Context, token string) error {
	if token == "" {
		return fmt.Errorf("telegram %w: token", util.ErrNotConfigured)
	}

	var err error
	t.Bot, err = bot.New(token, bot.WithSkipGetMe())
	if err != nil {
		return err
	}

	user, err := t.Bot.GetMe(ctx)
	if err != nil {
		return err
	}

	slog.Info("Connected to Telegram", "username", user.Username)
	return nil
}

// Send sends a message to the configured Telegram chat.
func (t *Telegram) Send(ctx context.Context, message string) error {
	if t.Bot == nil {
		return nil
	}

	_, err := t.Bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    t.ChatID,
		Text:      message,
		ParseMode: "markdown",
	})
	return err
}
