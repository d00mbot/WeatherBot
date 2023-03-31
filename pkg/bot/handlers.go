package bot

import (
	"context"
	"fmt"
	ct "subscription-bot/pkg/time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	commandStart          = "start"
	commandHelp           = "help"
	commandTime           = "time"
	locationButtonUnicode = 0x1F30E
	placeholderText       = "00 - 23"
)

var locationKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButtonLocation(string([]rune{locationButtonUnicode}) + "Send geolocation"),
	),
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	case commandHelp:
		return b.handleHelpCommand(message)
	case commandTime:
		return b.handleTimeCommand(message)
	default:
		b.handleUnknownCommand(message)
	}

	return nil
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	logger := b.container.GetLogger()

	msg := tgbotapi.NewMessage(message.Chat.ID, b.responses.Start)

	msg.ReplyMarkup = locationKeyboard

	if _, err := b.bot.Send(msg); err != nil {
		logger.Errorf("error sending message to telegram: %s ", err)
		return err
	}

	return nil
}

func (b *Bot) handleHelpCommand(message *tgbotapi.Message) error {
	logger := b.container.GetLogger()

	msg := tgbotapi.NewMessage(message.Chat.ID, b.responses.Help)

	if _, err := b.bot.Send(msg); err != nil {
		logger.Errorf("error sending message to telegram: %s ", err)
		return err
	}

	return nil
}

func (b *Bot) handleTimeCommand(message *tgbotapi.Message) error {
	logger := b.container.GetLogger()

	msg := tgbotapi.NewMessage(message.Chat.ID, b.responses.Time)
	msg.ReplyToMessageID = message.MessageID
	msg.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply:            true,
		InputFieldPlaceholder: placeholderText,
	}

	if _, err := b.bot.Send(msg); err != nil {
		logger.Errorf("error sending message to telegram: %s ", err)
		return err
	}

	return nil
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	logger := b.container.GetLogger()

	msg := tgbotapi.NewMessage(message.Chat.ID, b.responses.UnknownCommand)

	if _, err := b.bot.Send(msg); err != nil {
		logger.Errorf("error sending message to telegram: %s ", err)
		return err
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
		if err := b.handleOtherMessages(message); err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) handleLocationMessage(ctx context.Context, message *tgbotapi.Message) error {
	logger := b.container.GetLogger()

	msg := tgbotapi.NewMessage(message.Chat.ID, "")

	userExist, err := b.storage.CheckExist(ctx, b.client, message)
	if err != nil {
		return err
	}

	if !userExist {
		err := b.storage.Create(ctx, b.client, message, b.container.GetConfig())
		if err != nil {
			return err
		}
		msg.Text = b.responses.Location
	} else {
		err := b.storage.Update(ctx, b.client, message, b.container.GetConfig())
		if err != nil {
			return err
		}
		msg.Text = b.responses.LocationUpdate
	}

	if _, err := b.bot.Send(msg); err != nil {
		logger.Errorf("error sending message to telegram: %s ", err)
		return err
	}

	return nil
}

func (b *Bot) handleTimeMessage(ctx context.Context, message *tgbotapi.Message) error {
	logger := b.container.GetLogger()

	msg := tgbotapi.NewMessage(message.Chat.ID, "")

	userExist, err := b.storage.CheckExist(ctx, b.client, message)
	if err != nil {
		return err
	}

	if !userExist {
		msg.Text = b.responses.DefaultMessage
	} else {
		time := ct.CheckTime(message.Text)

		if time == "" {
			msg.Text = b.responses.WrongTime
		} else {
			msg.Text = fmt.Sprintf(b.responses.SuccessfulTime, time)

			if err := b.storage.UpdateTime(ctx, b.client, message, time); err != nil {
				return err
			}
		}
	}

	if _, err := b.bot.Send(msg); err != nil {
		logger.Errorf("error sending message to telegram: %s ", err)
		return err
	}

	return nil
}

func (b *Bot) handleOtherMessages(message *tgbotapi.Message) error {
	logger := b.container.GetLogger()

	msg := tgbotapi.NewMessage(message.Chat.ID, b.responses.DefaultMessage)

	if _, err := b.bot.Send(msg); err != nil {
		logger.Errorf("error sending message to telegram: %s ", err)
		return err
	}

	return nil
}
