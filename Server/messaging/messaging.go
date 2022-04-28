package messaging

import (
	"fmt"

	"api.mooody.me/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramMessaging struct {
	botApi       *tgbotapi.BotAPI
	targetChatId int64
	enabled      bool
}

func NewTelegramMessaging(enabled bool, token string, targetChatId int64) (*TelegramMessaging, error) {
	if enabled {
		return &TelegramMessaging{botApi: nil, targetChatId: targetChatId, enabled: enabled}, nil
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &TelegramMessaging{botApi: bot, targetChatId: targetChatId, enabled: enabled}, nil
}

func (m *TelegramMessaging) SendMessage(message string) {
	if !m.enabled {
		return
	}

	msg := tgbotapi.NewMessage(0, message)
	msg.ParseMode = "markdown"
	msg.ChatID = m.targetChatId

	_, err := m.botApi.Send(msg)
	if err != nil {
		fmt.Println(err)
	}
}

func (m *TelegramMessaging) SendNotification(event *models.Notification) {
	if !m.enabled {
		return
	}

	msg := tgbotapi.NewMessage(0, fmt.Sprintf("`%s` [%s]: %s", event.ChannelId, event.Title, event.Content))
	msg.ParseMode = "markdown"
	msg.ChatID = m.targetChatId

	_, err := m.botApi.Send(msg)
	if err != nil {
		fmt.Println(err)
	}
}
