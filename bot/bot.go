package bot

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"log"
	"time"
)

var bot *tgbotapi.BotAPI

var messagesChan chan tgbotapi.Chattable

func Run(mc chan tgbotapi.Chattable) {
	initBot()

	// bot 包使用messageChan进行消息的收发
	messagesChan = mc
	// debug日志开关
	bot.Debug = false

	go receiver()
	go sender()
}

// 初始化bot,失败会重试
func initBot() {
	err := errors.New("haven't init bot")
	for err != nil {
		bot, err = tgbotapi.NewBotAPI(viper.GetString("Token"))
		time.Sleep(5 * time.Second)
	}
	log.Println("init done")
}

func receiver() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Println(err)
		time.Sleep(5 * time.Second)
		go receiver()
		return
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		// 只回应来自管理员的消息
		if update.Message.Chat.ID != viper.GetInt64("AdminID") {
			continue
		}
		// 将消息交给命令路由
		go baseRouter(update)
	}
}

func sender() {
	for msg := range messagesChan {
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
}
