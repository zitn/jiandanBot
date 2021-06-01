package bot

import (
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"

	"jiandanBot/maker"
)

func baseRouter(update tgbotapi.Update) {
	// 只回应来自管理员的消息
	if update.Message.Chat.ID != viper.GetInt64("AdminID") {
		return
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
	// msg.ReplyToMessageID = update.Message.MessageID

	_, _ = botAPI.Send(msg)
}

func callbackRouter(update tgbotapi.Update) {
	commandAndData := strings.Fields(update.CallbackQuery.Data)
	switch commandAndData[0] {
	case "updateTucao":
		result := maker.UpdateTucao(update, commandAndData[1])

		// 返回提示
		_, _ = botAPI.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, result))
	default:
		_, _ = botAPI.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "command not found"))

	}
	return
}

func updateApiAddress(update tgbotapi.Update) {
	if update.Message.Text != "" {
		viper.Set("ApiAddress", update.Message.Text)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "api地址已更新为"+update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID
		_, _ = botAPI.Send(msg)
	}
}
