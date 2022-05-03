package messaging

import (
	"fmt"
	"log"
	"time"

	"api.mooody.me/common"
	"api.mooody.me/models/notifications"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramMessaging struct {
	botApi     *tgbotapi.BotAPI
	safeChatId int64
	safeUserId int64
	enabled    bool
}

func NewTelegramBot(enabled bool, token string, safeChatId int64, safeUserId int64) (*TelegramMessaging, error) {
	if !enabled {
		return &TelegramMessaging{botApi: nil, safeChatId: 0, safeUserId: 0, enabled: enabled}, nil
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Update bot commands each time
	mm := tgbotapi.NewSetMyCommands(tgbotapi.BotCommand{
		Command:     "ping",
		Description: "Ping!",
	}, tgbotapi.BotCommand{
		Command:     "version",
		Description: "Get the current version of the bot (git revison)",
	}, tgbotapi.BotCommand{
		Command:     "status",
		Description: "Get MoodyAPI status",
	})
	bot.Request(mm)

	return &TelegramMessaging{botApi: bot, safeChatId: safeChatId, safeUserId: safeUserId, enabled: enabled}, nil
}

func (m *TelegramMessaging) SendMessage(message string) {
	if !m.enabled {
		return
	}

	msg := tgbotapi.NewMessage(0, message)
	msg.ParseMode = "markdown"
	msg.ChatID = m.safeChatId

	_, err := m.botApi.Send(msg)
	if err != nil {
		fmt.Println(err)
	}
}

func (m *TelegramMessaging) SendNotification(event *notifications.Notification) {
	if !m.enabled {
		return
	}

	msg := tgbotapi.NewMessage(0, fmt.Sprintf("`%d` [%s]: %s", event.ChannelId, event.Title, event.Content))
	msg.ParseMode = "markdown"
	msg.ChatID = m.safeChatId

	_, err := m.botApi.Send(msg)
	if err != nil {
		fmt.Println(err)
	}
}

func (m *TelegramMessaging) HandleBotCommand() {
	if !m.enabled {
		return
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := m.botApi.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ReplyToMessageID = update.Message.MessageID

		if update.Message.Chat.ID != m.safeChatId && update.Message.Chat.ID != m.safeUserId {
			msg.Text = "This bot is only for Moody's chat group."
		} else {
			switch update.Message.Command() {
			case "ping":
				msg.Text = "🏓"
			case "version":
				msg.Text = "Server Version: " + common.ServerRevision
			case "status":
				msg.Text = fmt.Sprintf("API server has been running for %d minute(s)", int(time.Now().Sub(common.StartTime).Minutes()))
			default:
				msg.Text = "What are you talking about?"
			}
		}

		if _, err := m.botApi.Send(msg); err != nil {
			log.Fatal(err)
		}
	}
}
