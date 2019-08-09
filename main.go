package main

import (
	"myTeleBot/bot"
	"myTeleBot/crawler"
	"myTeleBot/maker"
	"sync"
)

func main() {
	go bot.Run()
	go crawler.GetJiandan()
	go maker.Jiandan()
	go maker.UpdateTucao()
	// 暴力的防止主进程退出
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
