package messaging

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"api.mooody.me/api"
	"api.mooody.me/common"
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
	if api.APIServer.LastLightState.On {
		api.APIServer.LastLightState.On = false
		api.APIServer.BroadcastLightState(api.APIServer.LastLightState)
		msg.Text = from + " 把灯关了"
	} else {
		msg.Text = "灯没开"
	}
}

func onLightOnAction(msg *tgbotapi.MessageConfig, from string) {
	if !api.APIServer.LastLightState.On {
		api.APIServer.LastLightState.On = true
		api.APIServer.BroadcastLightState(api.APIServer.LastLightState)
		msg.Text = from + " 把灯打开了"
	} else {
		msg.Text = "灯已经开着了"
	}
}

func onGetLightAction(msg *tgbotapi.MessageConfig) {
	if api.APIServer.LastLightState.On {
		msg.Text = "灯亮着"
		msg.Text += fmt.Sprintf(" \\(%d\\)", api.APIServer.LastLightState.Brightness)
	} else {
		msg.Text = "灯关着"
	}

	if api.APIServer.LastLightState.GetColored() != nil {
		color := api.APIServer.LastLightState.GetColored()
		msg.Text += fmt.Sprintf("\nColor: %d, %d, %d", color.Red, color.Green, color.Blue)
	} else {
		msg.Text += "\nColor: Warm White"
	}
}

func onColorAction(msg *tgbotapi.MessageConfig, from string, color []string) {
	if len(color) == 1 {
		if color[0] == "暖白" || color[0] == "warm" || color[0] == "ww" || color[0] == "warmwhite" {
			api.APIServer.LastLightState.Mode = &light.LightState_Warmwhite{Warmwhite: true}
			msg.Text = from + " 把灯调成了暖白"
		} else if color[0] == "red" || color[0] == "红" || color[0] == "r" || color[0] == "R" || color[0] == "红色" {
			api.APIServer.LastLightState.Mode = &light.LightState_Colored{Colored: &light.LightColor{Red: 255, Green: 0, Blue: 0}}
			msg.Text = from + " 把灯调成了红色"
		} else if color[0] == "green" || color[0] == "绿" || color[0] == "g" || color[0] == "G" || color[0] == "绿色" {
			api.APIServer.LastLightState.Mode = &light.LightState_Colored{Colored: &light.LightColor{Red: 0, Green: 255, Blue: 0}}
			msg.Text = from + " 把灯调成了绿色"
		} else if color[0] == "blue" || color[0] == "蓝" || color[0] == "b" || color[0] == "B" || color[0] == "蓝色" {
			api.APIServer.LastLightState.Mode = &light.LightState_Colored{Colored: &light.LightColor{Red: 0, Green: 0, Blue: 255}}
			msg.Text = from + " 把灯调成了蓝色"
		} else if color[0] == "yellow" || color[0] == "黄" || color[0] == "y" || color[0] == "Y" || color[0] == "黄色" {
			api.APIServer.LastLightState.Mode = &light.LightState_Colored{Colored: &light.LightColor{Red: 255, Green: 255, Blue: 0}}
			msg.Text = from + " 把灯调成了黄色"
		} else if color[0] == "cyan" || color[0] == "青" || color[0] == "c" || color[0] == "C" || color[0] == "青色" {
			api.APIServer.LastLightState.Mode = &light.LightState_Colored{Colored: &light.LightColor{Red: 0, Green: 255, Blue: 255}}
			msg.Text = from + " 把灯调成了青色"
		} else if color[0] == "purple" || color[0] == "紫" || color[0] == "p" || color[0] == "P" || color[0] == "紫色" {
			api.APIServer.LastLightState.Mode = &light.LightState_Colored{Colored: &light.LightColor{Red: 255, Green: 0, Blue: 255}}
			msg.Text = from + " 把灯调成了紫色"
		} else if color[0] == "black" || color[0] == "黑" || color[0] == "k" || color[0] == "K" || color[0] == "黑色" {
			api.APIServer.LastLightState.Mode = &light.LightState_Colored{Colored: &light.LightColor{Red: 0, Green: 0, Blue: 0}}
			msg.Text = from + " 把灯调成了黑色（？）"
		} else if color[0] == "white" || color[0] == "白" || color[0] == "w" || color[0] == "W" || color[0] == "白色" {
			api.APIServer.LastLightState.Mode = &light.LightState_Colored{Colored: &light.LightColor{Red: 255, Green: 255, Blue: 255}}
			msg.Text = from + " 把灯调成了白色"
		} else if color[0] == "色" {
			msg.Text = "好！"
		} else if strings.HasPrefix(color[0], "#") {
			c, err := common.ParseHexColorFast(color[0])
			if err != nil {
				msg.Text = "颜色格式不对：" + err.Error()
			} else {
				api.APIServer.LastLightState.Mode = &light.LightState_Colored{Colored: &light.LightColor{Red: uint32(c.R), Green: uint32(c.G), Blue: uint32(c.B)}}
				msg.Text = from + " 把灯调成了 \\" + color[0]
			}
		} else {
			msg.Text = "不对"
			return
		}
	} else if len(color) != 3 {
		msg.Text = "不够色"
		return
	} else {
		red, err := strconv.Atoi(color[0])
		if err != nil || red < 0 || red > 255 {
			msg.Text = "红色不对: " + color[0]
			return
		}

		green, err := strconv.Atoi(color[1])
		if err != nil || green < 0 || green > 255 {
			msg.Text = "绿色不对: " + color[1]
			return
		}

		blue, err := strconv.Atoi(color[2])
		if err != nil || blue < 0 || blue > 255 {
			msg.Text = "蓝色不对: " + color[2]
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
