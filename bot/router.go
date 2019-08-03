package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"myTeleBot/maker"
	"strings"
)

func baseRouter(update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		commandAndData := strings.Fields(update.CallbackQuery.Data)
		switch commandAndData[0] {
		case "updateTucao":
			err := maker.UpdateTucao(viper.GetInt64("ChannelID"), update.CallbackQuery.Message.MessageID, commandAndData[1], update.CallbackQuery.Message.Caption)
			if err != nil {
				// 返回错误
				_, _ = bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, err.Error()))

			}
		default:
			_, _ = bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "command not found"))

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

	_, _ = bot.Send(msg)
}

func updateApiAddress(update tgbotapi.Update) {
	viper.Set("ApiAddress", update.Message.Text)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "api地址已更新")
	msg.ReplyToMessageID = update.Message.MessageID
	_, _ = bot.Send(msg)
}
