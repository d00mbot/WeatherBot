package bot

import (
	"context"
	"subscription-bot/internal/time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	commandStart = "start"
	commandHelp  = "help"
	commandTime  = "time"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.responseStartCommand(message)
	case commandHelp:
		return b.responseHelpCommand(message)
	case commandTime:
		return b.responseTimeCommand(message)
	default:
		b.responseUnknownCommand(message)
	}

	return nil
}

func (b *Bot) handleMessage(ctx context.Context, message *tgbotapi.Message) error {
	if message.Location != nil {
		if err := b.handleLocationMessage(ctx, message); err != nil {
			return err
		}
	}

	if message.ReplyToMessage != nil {
		if err := b.handleTimeMessage(ctx, message); err != nil {
			return err
		}
	}

	if message.Location == nil && message.ReplyToMessage == nil {
		if err := b.responseDefaultMessage(message); err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) handleLocationMessage(ctx context.Context, message *tgbotapi.Message) error {
	userExist, err := b.storage.CheckExist(ctx, b.client, message)
	if err != nil {
		return err
	}

	if !userExist {
		if err := b.storage.Create(ctx, b.client, message); err != nil {
			return err
		}
	} else {
		if err := b.storage.Update(ctx, b.client, message); err != nil {
			return err
		}
	}

	if err := b.responseLocationMessage(userExist, message); err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleTimeMessage(ctx context.Context, message *tgbotapi.Message) error {
	userExist, err := b.storage.CheckExist(ctx, b.client, message)
	if err != nil {
		return err
	}

	if !userExist {
		if err := b.responseDefaultMessage(message); err != nil {
			return err
		}
	} else {
		time := time.CheckTime(message.Text)

		if time != "" {
			if err := b.storage.UpdateTime(ctx, b.client, message, time); err != nil {
				return err
			}
		}

		if err := b.responseTimeMessage(time, message); err != nil {
			return err
		}
	}

	return nil
}
