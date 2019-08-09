package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"myTeleBot/channel"
	"strings"
)

func baseRouter(update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		commandAndData := strings.Fields(update.CallbackQuery.Data)
		switch commandAndData[0] {
		case "updateTucao":
			// 返回提示
			channel.RequireUpdateTucaoChannel <- commandAndData[1]
			_, _ = botAPI.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "正在更新，请勿重复点击"))
		default:
			_, _ = botAPI.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "command not found"))

		}
	}

	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "updateApi":
			updateApiAddress(update)
			return
		}

		// todo 处理其他命令,以及寻求更优雅的实现
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	//msg.ReplyToMessageID = update.Message.MessageID

	_, _ = botAPI.Send(msg)
}

func updateApiAddress(update tgbotapi.Update) {
	if update.Message.Text != "" {
		viper.Set("ApiAddress", update.Message.Text)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "api地址已更新为"+update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID
		_, _ = botAPI.Send(msg)
	}

}
