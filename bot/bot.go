package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"log"
	"time"
)

var (
	// 包级别变量，方便包内调用
	botAPI = initBot()
)

// 初始化bot
func initBot() *tgbotapi.BotAPI {
	api, err := tgbotapi.NewBotAPI(viper.GetString("Token"))
	if err != nil {
		log.Panic(err)
	}
	log.Println("init done")
	return api
}

// 初始化 bot 的各个服务
func Run() {
	// debug日志开关
	botAPI.Debug = false
	go receiver()
	go sender()
	go commentSender()
}

func receiver() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := botAPI.GetUpdatesChan(u)

	if err != nil {
		log.Println(err)
		time.Sleep(5 * time.Second)
		go receiver()
		return
	}

	for update := range updates {
		if update.CallbackQuery != nil {
			go callbackRouter(update)
		}
		if update.Message == nil {
			continue
		}
		// 将消息交给路由,处理下一条消息
		go baseRouter(update)
	}
}
