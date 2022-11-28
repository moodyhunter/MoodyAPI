package messaging

import (
	"context"
	"fmt"
	"strconv"

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

func onLightOffAction(msg *tgbotapi.MessageConfig, from string) {
	api.APIServer.LastLightState.On = false
	api.APIServer.BroadcastLightState(api.APIServer.LastLightState)
	msg.Text = from + " 把灯关了"
}

func onLightOnAction(msg *tgbotapi.MessageConfig, from string) {
	api.APIServer.LastLightState.On = true
	api.APIServer.BroadcastLightState(api.APIServer.LastLightState)
	msg.Text = from + " 把灯打开了"
}

func onGetLightAction(msg *tgbotapi.MessageConfig) {
	if api.APIServer.LastLightState.On {
		msg.Text = "灯亮着"
	} else {
		msg.Text = "灯关着"
	}

	msg.Text += fmt.Sprintf(" \\(brightness: %d\\)", api.APIServer.LastLightState.Brightness)

	if api.APIServer.LastLightState.GetColored() != nil {
		color := api.APIServer.LastLightState.GetColored()
		msg.Text += fmt.Sprintf("\nColor: %d, %d, %d", color.Red, color.Green, color.Blue)
	} else {
		msg.Text += "\nColor: Warm White"
	}
}

func onColorAction(msg *tgbotapi.MessageConfig, from string, color []string) {
	if len(color) == 1 {
		if color[0] == "白" || color[0] == "warm" || color[0] == "ww" || color[0] == "warmwhite" || color[0] == "暖白" {
			api.APIServer.LastLightState.Mode = &light.LightState_Warmwhite{Warmwhite: true}
			msg.Text = from + " 把灯调成了暖白"
		} else {
			msg.Text = color[0] + " 不对吧？"
			return
		}
	} else if len(color) != 3 {
		msg.Text = "不够色"
		return
	} else {
		red, err := strconv.Atoi(color[0])
		if err != nil || red < 0 || red > 255 {
			msg.Text = "红色不对"
			return
		}

		green, err := strconv.Atoi(color[1])
		if err != nil || green < 0 || green > 255 {
			msg.Text = "绿色不对"
			return
		}

		blue, err := strconv.Atoi(color[2])
		if err != nil || blue < 0 || blue > 255 {
			msg.Text = "蓝色不对"
			return
		}

		api.APIServer.LastLightState.Mode = &light.LightState_Colored{
			Colored: &light.LightColor{
				Red:   uint32(red),
				Green: uint32(green),
				Blue:  uint32(blue),
			},
		}
		msg.Text = from + " 把灯调成了 " + color[0] + ", " + color[1] + ", " + color[2]
	}

	api.APIServer.BroadcastLightState(api.APIServer.LastLightState)
}
