package bot

import (
	"log"
	"net/http"
	"net/url"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	// 包级别变量，方便包内调用
	botAPI *tgbotapi.BotAPI
)

// 初始化bot
func initBot() *tgbotapi.BotAPI {
	if viper.GetString("telegram_proxy") == "" {
		api, err := tgbotapi.NewBotAPI(viper.GetString("Token"))
		if err != nil {
			log.Panic(err)
		}
		logrus.Info("init done, start working")
		return api
	} else {
		proxy := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(viper.GetString("telegram_proxy"))
		}
		httpClient := &http.Client{
			Transport: &http.Transport{
				Proxy: proxy,
			},
		}
		api, err := tgbotapi.NewBotAPIWithClient(viper.GetString("Token"), httpClient)
		if err != nil {
			log.Panic(err)
		}
		return api
	}
}

// 初始化 bot 的各个服务
func Run() {
	botAPI = initBot()
	// debug日志开关
	botAPI.Debug = false
	go sender()
	go receiver()
	go commentSender()
}

func receiver() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := botAPI.GetUpdatesChan(u)

	if err != nil {
		logrus.WithField("func", "receiver").WithField("err in", "GetUpdatesChan").Panicln(err)
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
