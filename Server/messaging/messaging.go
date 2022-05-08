package messaging

import (
	"context"
	"fmt"
	"log"
	"time"

	"api.mooody.me/common"
	"api.mooody.me/db"
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
		Command:     "status",
		Description: "Get MoodyAPI status",
	}, tgbotapi.BotCommand{
		Command:     "channels",
		Description: "List notification channels",
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

	// get channel name from event's channel Id
	channelName := "<unknown>"
	channel, err := db.GetNotificationChannelById(context.Background(), event.ChannelId)
	if err != nil {
		fmt.Println(err)
	}
	channelName = channel.Name

	msg := tgbotapi.NewMessage(0, fmt.Sprintf("*New Message From Channel \"%s\"*\n\n*Title:* %s\n*Content:* %s", channelName, event.Title, event.Content))
	msg.ParseMode = "markdown"
	msg.ChatID = m.safeChatId

	_, err = m.botApi.Send(msg)
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
		msg.Text = "`" + update.Message.Command() + "` is not implemented yet"
		msg.ParseMode = "markdownv2"

		if update.Message.Chat.ID != m.safeChatId && update.Message.Chat.ID != m.safeUserId {
			msg.Text = "This bot is only for Moody's chat group."
		} else {
			switch update.Message.Command() {
			case "ping":
				msg.Text = "🏓"
			case "status":
				msg.Text = "Server Revision: `" + common.ServerRevision + "`\n"
				msg.Text += fmt.Sprintf("Uptime: `%d` minute\\(s\\)", int(time.Now().Sub(common.StartTime).Minutes()))
			case "channels":
				onChannelsAction(&msg)
			default:
				msg.Text = "What are you talking about?"
			}
		}

		if _, err := m.botApi.Send(msg); err != nil {
			log.Println(err)
		}
	}
}
