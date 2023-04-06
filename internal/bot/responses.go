package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	locationButtonUnicode = 0x1F30E
	locationButtonText    = "Send geolocation"
	placeholderText       = "00 - 23"
)

var locationKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButtonLocation(string([]rune{locationButtonUnicode}) + locationButtonText),
	),
)

func (b *Bot) responseStartCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.responses.Start)

	msg.ReplyMarkup = locationKeyboard

	_, err := b.bot.Send(msg)
	if err != nil {
		b.container.GetLogger().Errorf("faild to send message to telegram: %s ", err)
		return err
	}

	return nil
}

func (b *Bot) responseHelpCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.responses.Help)

	_, err := b.bot.Send(msg)
	if err != nil {
		b.container.GetLogger().Errorf("faild to send message to telegram: %s ", err)
		return err
	}

	return nil
}

func (b *Bot) responseTimeCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.responses.Time)
	msg.ReplyToMessageID = message.MessageID
	msg.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply:            true,
		InputFieldPlaceholder: placeholderText,
	}

	_, err := b.bot.Send(msg)
	if err != nil {
		b.container.GetLogger().Errorf("faild to send message to telegram: %s ", err)
		return err
	}

	return nil
}

func (b *Bot) responseUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.responses.UnknownCommand)

	_, err := b.bot.Send(msg)
	if err != nil {
		b.container.GetLogger().Errorf("faild to send message to telegram: %s ", err)
		return err
	}

	return nil
}

func (b *Bot) responseDefaultMessage(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.responses.DefaultMessage)

	_, err := b.bot.Send(msg)
	if err != nil {
		b.container.GetLogger().Errorf("faild to send message to telegram: %s ", err)
		return err
	}

	return nil
}

func (b *Bot) responseLocationMessage(userExist bool, message *tgbotapi.Message) error {
	if !userExist {
		msg := tgbotapi.NewMessage(message.Chat.ID, b.responses.UserCreated)

		_, err := b.bot.Send(msg)
		if err != nil {
			b.container.GetLogger().Errorf("faild to send message to telegram: %s ", err)
			return err
		}
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, b.responses.UserUpdated)

		_, err := b.bot.Send(msg)
		if err != nil {
			b.container.GetLogger().Errorf("faild to send message to telegram: %s ", err)
			return err
		}
	}

	return nil
}

func (b *Bot) responseTimeMessage(time string, message *tgbotapi.Message) error {
	if time == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, b.responses.WrongTime)
		msg.ReplyToMessageID = message.MessageID
		msg.ReplyMarkup = tgbotapi.ForceReply{
			ForceReply:            true,
			InputFieldPlaceholder: placeholderText,
		}

		_, err := b.bot.Send(msg)
		if err != nil {
			b.container.GetLogger().Errorf("faild to send message to telegram: %s ", err)
			return err
		}
	} else {
		msgText := fmt.Sprintf(b.responses.TimeUpdated, time)
		msg := tgbotapi.NewMessage(message.Chat.ID, msgText)

		_, err := b.bot.Send(msg)
		if err != nil {
			b.container.GetLogger().Errorf("faild to send message to telegram: %s ", err)
			return err
		}
	}

	return nil
}
