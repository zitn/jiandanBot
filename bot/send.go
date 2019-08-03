package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"myTeleBot/types"
)

// 普通消息sender
func sender(messagesChan chan tgbotapi.Chattable) {
	for msg := range messagesChan {
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
}

// 煎蛋帖子sender
func CommentSender(commentChannel chan types.CommentMessage) {
	for message := range commentChannel {
		CommentResponse, err := bot.Send(message.CommentMessage)
		if err != nil {
			log.Println(err)
		}
		message.TucaoMessage.ReplyToMessageID = CommentResponse.MessageID
		TucaoResponse, err := bot.Send(message.TucaoMessage)
		if err != nil {
			log.Println(err)
		}
		if !message.HaveTucao {
			// todo 十分钟之后主动更新
		}
	}
}

// 更新煎蛋吐槽
func tucaoUpdater() {

}
