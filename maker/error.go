package maker

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"myTeleBot/channel"
)

func ErrorMaker() {
	for errorMessage := range channel.ErrorMessage {
		newErrorMessage := tgbotapi.NewMessage(viper.GetInt64("150606003"), errorMessage.Error())
		channel.NormalMessageChannel <- newErrorMessage
	}
}
