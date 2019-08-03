package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"myTeleBot/bot"
	"myTeleBot/crawler"
	"myTeleBot/maker"
	"myTeleBot/types"
	"sync"
)

func main() {
	// 所有处理好待发送的消息均放入此channel中
	messageChan := make(chan tgbotapi.Chattable, 100)

	// 煎蛋帖子channel
	commentMessageChannel := make(chan types.CommentMessage, 40)

	// 煎蛋更新需求channel
	// todo

	// 启动bot,监听消息
	bot.Run(messageChan, commentMessageChannel)

	// todo 优化逻辑
	// 煎蛋爬虫channel
	commentsChan := make(chan types.Comment, 30)

	// 启动煎蛋爬虫
	go crawler.GetJiandan(commentsChan)

	// 启动煎蛋maker
	go maker.Jiandan(messageChan, commentsChan)

	// todo 测试
	//bot.SendTester()

	// 防止主进程退出
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
