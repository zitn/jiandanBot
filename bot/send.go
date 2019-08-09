package bot

import (
	"log"
	"myTeleBot/channel"
)

// 普通消息sender
func sender() {
	for msg := range channel.NormalMessageChannel {
		_, err := botAPI.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
}

// 煎蛋帖子sender
func CommentSender() {
	for message := range channel.CommentMessageChannel {
		CommentResponse, err := botAPI.Send(message.CommentMessage)
		if err != nil {
			log.Println(err)
		}
		message.TucaoMessage.ReplyToMessageID = CommentResponse.MessageID
		_, err = botAPI.Send(message.TucaoMessage)
		if err != nil {
			log.Println(err)
		}
	}
}
