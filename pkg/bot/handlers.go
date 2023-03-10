package bot

import (
	"context"
	"fmt"
	database "subscription-bot/pkg/database"
	ct "subscription-bot/pkg/time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

const (
	commandStart = "start"
	commandTime  = "time"
)

var locationKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButtonLocation(string([]rune{0x1F30E}) + "Send geolocation"),
	),
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	case commandTime:
		return b.handleTimeCommand(message)
	default:
		b.handleUnknownCommand(message)
	}

	return nil
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.responses.Start)

	msg.ReplyMarkup = locationKeyboard

	if _, err := b.bot.Send(msg); err != nil {
		log.Errorf("error sending message to telegram: %s ", err)
		return err
	}

	return nil
}

func (b *Bot) handleTimeCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.responses.Time)
	msg.ReplyToMessageID = message.MessageID
	msg.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply:            true,
		InputFieldPlaceholder: "00 - 23",
	}

	if _, err := b.bot.Send(msg); err != nil {
		log.Errorf("error sending message to telegram: %s ", err)
		return err
	}

	return nil
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.responses.UnknownCommand)

	if _, err := b.bot.Send(msg); err != nil {
		log.Errorf("error sending message to telegram: %s ", err)
		return err
	}

	return nil
}

func (b *Bot) handleOtherMessage(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.responses.OtherMessage)

	if _, err := b.bot.Send(msg); err != nil {
		log.Errorf("error sending message to telegram: %s ", err)
		return err
	}

	return nil
}

func (b *Bot) handleLocationMessage(ctx context.Context, message *tgbotapi.Message) error {
	if err := database.InsertSubscriber(ctx, b.client, message, b.weatherToken); err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, b.responses.Location)

	if _, err := b.bot.Send(msg); err != nil {
		log.Errorf("error sending message to telegram: %s ", err)
		return err
	}

	return nil
}

func (b *Bot) handleReplyMessage(ctx context.Context, message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "")

	time := ct.CheckTime(message.Text)

	if time == "" {
		msg.Text = b.responses.WrongTime
	} else {
		msg.Text = fmt.Sprintf(b.responses.SuccessfulTime, time)

		if err := database.UpdateUserTime(ctx, b.client, message, time); err != nil {
			return err
		}
	}

	if _, err := b.bot.Send(msg); err != nil {
		log.Errorf("error sending message to telegram: %s ", err)
		return err
	}

	return nil
}
