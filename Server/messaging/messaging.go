package messaging

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"api.mooody.me/common"
	"api.mooody.me/db"
	"api.mooody.me/models/notifications"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	botApi     *tgbotapi.BotAPI
	safeChatId int64
	safeUserId int64
}

func NewTelegramBot(token string, safeChatId int64, safeUserId int64) *TelegramBot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

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

	return &TelegramBot{botApi: bot, safeChatId: safeChatId, safeUserId: safeUserId}
}

func (bot *TelegramBot) Close() {
	bot.botApi.StopReceivingUpdates()
	log.Println("Telegram bot is closed.")
}

func (m *TelegramBot) SendMessage(message string) {
	msg := tgbotapi.NewMessage(0, message)
	msg.ParseMode = "markdown"
	msg.ChatID = m.safeChatId

	_, err := m.botApi.Send(msg)
	if err != nil {
		fmt.Println(err)
	}
}

func (m *TelegramBot) SendNotification(event *notifications.Notification) {
	// get channel name from event's channel Id
	channelName := "INVALID CHANNEL: " + strconv.FormatInt(event.ChannelId, 10)
	channel, err := db.GetNotificationChannelById(context.Background(), event.ChannelId)
	if err != nil {
		fmt.Println(err)
	} else {
		channelName = channel.Name
	}

	channelName = tgbotapi.EscapeText(tgbotapi.ModeMarkdown, channelName)
	event.Title = tgbotapi.EscapeText(tgbotapi.ModeMarkdown, event.Title)
	event.Content = tgbotapi.EscapeText(tgbotapi.ModeMarkdown, event.Content)

	msg := tgbotapi.NewMessage(0, fmt.Sprintf("*New Message From Channel \"%s\"*\n\n*Title:* %s\n*Content:* %s", channelName, event.Title, event.Content))
	msg.ParseMode = "markdown"

	if event.Private {
		msg.ChatID = m.safeUserId
	} else {
		msg.ChatID = m.safeChatId
	}

	_, err = m.botApi.Send(msg)
	if err != nil {
		fmt.Println(err)
	}
}

func (m *TelegramBot) ServeBotCommand() {
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
				msg.Text += fmt.Sprintf("Uptime: `%d` minute\\(s\\)", int(time.Since(common.StartTime).Minutes()))
			case "channels":
				onChannelsAction(&msg)
			default:
				continue
			}
		}

		if _, err := m.botApi.Send(msg); err != nil {
			log.Println(err)
		}
	}
}
