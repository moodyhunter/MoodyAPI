package messaging

import (
	"context"
	"fmt"

	"api.mooody.me/api"
	"api.mooody.me/db"
	"api.mooody.me/models/light"
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

func onLightOffAction(msg *tgbotapi.MessageConfig) {
	state := &light.LightState{
		On: false,
	}
	api.APIServer.BroadcastLightStatus(state)
}
