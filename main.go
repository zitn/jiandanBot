package main

import (
	"flag"
	"sync"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"jiandanBot/bot"
	"jiandanBot/channel"
	"jiandanBot/crawler"
	"jiandanBot/maker"
)

func init() {
	viper.SetConfigName("config") // 配置文件名
	viper.AddConfigPath("config") // 配置文件所在的路径
	viper.SetConfigType("json")   // 配置文件类型
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Panic(err)
	}
	debug := flag.Bool("debug", false, "debug")
	flag.Parse()
	if !*debug {
		errorMsg := new(ErrorMsg)
		logrus.SetOutput(errorMsg)
	}

}

func main() {
	go bot.Run()
	go crawler.GetJianDan()
	go maker.Jiandan()
	// 暴力的防止主进程退出
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

type ErrorMsg struct {
}

func (e *ErrorMsg) Write(p []byte) (n int, err error) {
	newErrorMessage := tgbotapi.NewMessage(viper.GetInt64("AdminID"), string(p))
	channel.NormalMessageChannel <- newErrorMessage
	return len(p), nil
}
