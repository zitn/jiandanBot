package bot

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"log"
	"myTeleBot/types"
	"time"
)

var (
	// 方便包内调用
	bot *tgbotapi.BotAPI
)

func Run(messagesChan chan tgbotapi.Chattable, commentChannel chan types.CommentMessage) {
	initBot()

	// debug日志开关
	bot.Debug = true

	go receiver()
	go sender(messagesChan)
	go CommentSender(commentChannel)
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
		// 将消息交给路由,处理下一条消息
		go baseRouter(update)
	}
}
