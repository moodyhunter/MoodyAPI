package messaging

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
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

var NonCommandVerbs = []string{
	"灯",
	"开灯",
	"关灯",
	"色", "color",
}

func NewTelegramBot(token string, safeChatId int64, safeUserId int64) *TelegramBot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	mm := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{Command: "ping", Description: "Ping!"},
		tgbotapi.BotCommand{Command: "status", Description: "Get MoodyAPI status"},
		tgbotapi.BotCommand{Command: "channels", Description: "List notification channels"},
		tgbotapi.BotCommand{Command: "light_off", Description: "Turn off the light"},
		tgbotapi.BotCommand{Command: "light_on", Description: "Turn on the light"},
		tgbotapi.BotCommand{Command: "get_light", Description: "Get light status"},
		tgbotapi.BotCommand{Command: "pin", Description: "Pin a message"},
	)
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

	msg := tgbotapi.NewMessage(0, fmt.Sprintf("*<%s>* - %s\n %s", channelName, event.Title, event.Content))
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
		if update.Message == nil {
			continue
		}

		command := ""
		args := []string{} // only supported for NonCommandVerbs

		if update.Message.IsCommand() {
			command = update.Message.Command()
		} else {
			tmp_args := strings.Split(update.Message.Text, " ")
			if len(tmp_args) <= 0 {
				continue
			}

			tmp_command := tmp_args[0]

			if !strings.HasPrefix(tmp_command, "/") {
				continue
			}

			tmp_command = strings.TrimPrefix(tmp_command, "/")

			for _, verb := range NonCommandVerbs {
				if tmp_command == verb {
					command = tmp_command
					args = tmp_args[1:]
					break
				}
			}
		}

		if command == "" {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ReplyToMessageID = update.Message.MessageID
		msg.Text = "`" + command + "` 是什么？"
		msg.ParseMode = "markdownv2"

		if update.Message.Chat.ID != m.safeChatId && update.Message.Chat.ID != m.safeUserId {
			msg.Text = "This bot is only for Moody's chat group."
		} else {
			switch command {
			case "ping":
				msg.Text = "🏓"
			case "status":
				msg.Text = "Server Revision: `" + common.ServerRevision + "`\n"
				msg.Text += fmt.Sprintf("Uptime: `%d` minute\\(s\\)", int(time.Since(common.StartTime).Minutes()))
			case "channels":
				onChannelsAction(&msg)
			case "light_off", "关灯":
				onLightOffAction(&msg)
			case "light_on", "开灯":
				onLightOnAction(&msg)
			case "get_light", "灯":
				onGetLightAction(&msg)
			case "色", "color":
				onColorAction(&msg, args)
			case "pin":
				if update.Message.ReplyToMessage != nil {
					_, err := m.botApi.Request(tgbotapi.PinChatMessageConfig{
						ChatID:              update.Message.Chat.ID,
						MessageID:           update.Message.ReplyToMessage.MessageID,
						DisableNotification: true,
					})
					if err != nil {
						msg.Text = err.Error()
					} else {
						msg.Text = "好！"
					}
				} else {
					msg.Text = "坏！"
				}
			default:
				continue
			}
		}

		msg.Text = tgbotapi.EscapeText(tgbotapi.ModeMarkdown, msg.Text)

		if _, err := m.botApi.Send(msg); err != nil {
			log.Println(err)
		}
	}
}
