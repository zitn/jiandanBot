package bot

import (
	"log"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"jiandanBot/channel"
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
func commentSender() {
	log1 := logrus.WithField("func", "commentSender")
	for message := range channel.CommentMessageChannel {
		CommentResponse, err := botAPI.Send(message.CommentMessage)
		if err != nil {
			if strings.Contains(err.Error(), "Too Many") {
				time.Sleep(1 * time.Second)
				channel.CommentMessageChannel <- message
				continue
			} else {
				log1.WithField("err in", "botAPI.Send CommentMessage").WithField("message is", message).Error(err)
				continue
			}
		}
		message.TucaoMessage.ReplyToMessageID = CommentResponse.MessageID
		_, err = botAPI.Send(message.TucaoMessage)
		if err != nil {
			if strings.Contains(err.Error(), "Too Many") {
				time.Sleep(10 * time.Second)
				continue
			}
			log1.WithField("err in", "botAPI.Send TucaoMessage").WithField("message is", message).Error(err)
		}
		// 睡5s防止发送过快
		time.Sleep(5 * time.Second)

	}
}
