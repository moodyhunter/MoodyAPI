package messaging

import (
	"context"
	"fmt"

	"api.mooody.me/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func onChannelsAction(msg *tgbotapi.MessageConfig) {
	channels, err := db.ListNotificationChannels(context.Background())

	if err != nil {
		msg.Text = fmt.Sprintf("Error: %s", err.Error())
		return
	}

	msg.Text = ""
	for _, channel := range channels {
		msg.Text += fmt.Sprintf("`%d` \\- `%s`\n", channel.Id, channel.Name)
	}
}
