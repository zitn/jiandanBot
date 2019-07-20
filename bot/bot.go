package bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"log"
	"time"
)

var bot *tgbotapi.BotAPI

func Run() {
	initBot()

	// todo debug日志开关
	bot.Debug = false

	go listenUpdate()
}

func initBot() {
	var err error
	bot, err = tgbotapi.NewBotAPI(viper.GetString("Token"))
	if err != nil {
		log.Panic(err)
	}

	log.Println("init done")
}

func listenUpdate() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Println(err)
		time.Sleep(5 * time.Second)
		listenUpdate()
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "updateApi":
				viper.Set("ApiAddress", update.Message.Text)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "api地址已更新")
				msg.ReplyToMessageID = update.Message.MessageID
				_, _ = bot.Send(msg)
				continue
			}

			// todo 处理命令消息
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID

		_, _ = bot.Send(msg)
	}
}

func Sender(messagesChan chan tgbotapi.Chattable) {
	for msg := range messagesChan {
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
}

//func SendTester() {
//	var testps []interface{}
//	testP := tgbotapi.NewInputMediaVideo("https://wx4.sinaimg.cn/mw1024/dc106893ly1g536eh99l0g20g00904qx.gif")
//	testP2 := tgbotapi.NewInputMediaVideo("https://wx4.sinaimg.cn/mw1024/dc106893ly1g536eh99l0g20g00904qx.gif")
//	testP3 := tgbotapi.NewInputMediaVideo("https://wx4.sinaimg.cn/mw1024/dc106893ly1g536eh99l0g20g00904qx.gif")
//
//	testps = append(testps, testP)
//	testps = append(testps, testP2)
//	testps = append(testps, testP3)
//	testMsg := tgbotapi.NewMediaGroup(viper.GetInt64("AdminID"), testps)
//	_, _ = bot.Send(testMsg)
//}
