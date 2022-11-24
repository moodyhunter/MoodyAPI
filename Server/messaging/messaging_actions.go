package messaging

import (
	"context"
	"fmt"

	"api.mooody.me/api"
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

func onLightOffAction(msg *tgbotapi.MessageConfig) {
	api.APIServer.LastLightState.On = false
	api.APIServer.BroadcastLightState(api.APIServer.LastLightState)
	msg.Text = "Light is off"
}

func onLightOnAction(msg *tgbotapi.MessageConfig) {
	api.APIServer.LastLightState.On = true
	api.APIServer.BroadcastLightState(api.APIServer.LastLightState)
	msg.Text = "Light is on"
}

func onGetLightAction(msg *tgbotapi.MessageConfig) {
	if api.APIServer.LastLightState.On {
		msg.Text = "Light is on"
	} else {
		msg.Text = "Light is off"
	}

	msg.Text += fmt.Sprintf(" (brightness: %d)", api.APIServer.LastLightState.Brightness)

	if api.APIServer.LastLightState.GetColored() != nil {
		color := api.APIServer.LastLightState.GetColored()
		msg.Text += fmt.Sprintf("\nColor: %d, %d, %d", color.Red, color.Green, color.Blue)
	} else {
		msg.Text += "\nColor: Warm White"
	}
}
