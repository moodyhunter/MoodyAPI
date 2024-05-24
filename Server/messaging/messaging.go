package messaging

import (
	"context"
	"fmt"
	"log"
	"slices"
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

var six = []string{
	"6",
	"６",
	"6️⃣",
	"六",
	"陆",
	"陸",
	"⁶",
	"₆",
	"Ⅵ", "ⅵ",
	"⑥", "❻", "➅", "➏",
	"⑹",
	"Ƅ",
	"㊅", "㈥",
	"𐄞", "𐄧", "𐄰", "𐄌", "𒐨",
	"𝟔", "𝟞", "𝟨", "𝟲", "𝟼", "🀕",
	"🃖", "🂶", "🂦", "🃆", "🯶",
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
	event.Content = tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, event.Content)

	msg := tgbotapi.NewMessage(0, fmt.Sprintf("*\\<%s\\>* \\- %s\n%s", channelName, event.Title, event.Content))
	msg.ParseMode = tgbotapi.ModeMarkdownV2

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

		if update.Message.IsCommand() {
			command = update.Message.Command()
		} else if update.Message.Sticker != nil {
			if slices.Contains(six, update.Message.Sticker.Emoji) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				msg.ReplyToMessageID = update.Message.MessageID
				msg.Text = "单走一个 6，傻逼。"
				m.botApi.Send(msg)
			}
		} else {
			tmp_args := strings.Split(update.Message.Text, " ")
			if len(tmp_args) <= 0 {
				continue
			}

			tmp_command := tmp_args[0]

			if !strings.HasPrefix(tmp_command, "/") {
				if slices.Contains(six, tmp_command) {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
					msg.ReplyToMessageID = update.Message.MessageID
					msg.Text = "单走一个 6，傻逼。"
					m.botApi.Send(msg)
				}
				continue
			}

			tmp_command = strings.TrimPrefix(tmp_command, "/")

			for _, verb := range NonCommandVerbs {
				if tmp_command == verb {
					command = tmp_command
					break
				}
			}
		}

		if command == "" {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ReplyToMessageID = update.Message.MessageID
		msg.Text = "不认识 `" + command + "`"
		msg.ParseMode = tgbotapi.ModeMarkdownV2

		if update.Message.Chat.ID != m.safeChatId && update.Message.Chat.ID != m.safeUserId {
			msg.Text = "This bot is only for Moody's chat group."
		} else {
			switch command {
			case "ping":
				msg.Text = "🏓"
			case "status":
				msg.Text = "`" + common.ServerRevision + "`\n"
				msg.Text += fmt.Sprintf("`%d` 分钟了", int(time.Since(common.StartTime).Minutes()))
			case "channels":
				onChannelsAction(&msg)
			case "light_off", "关灯", "light_on", "开灯", "get_light", "灯":
				msg.Text = "没有灯了捏"
			case "色", "color":
				msg.Text = "不能色了喔"
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
						msg.Text = "" // empty message
					}
				} else {
					msg.Text = "不行"
				}
			default:
				continue
			}
		}

		if msg.Text != "" {
			msg.Text = tgbotapi.EscapeText(tgbotapi.ModeMarkdown, msg.Text)

			if _, err := m.botApi.Send(msg); err != nil {
				log.Println(err)
			}
		}
	}
}
