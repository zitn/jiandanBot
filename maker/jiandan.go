package maker

import (
	"bytes"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"myTeleBot/types"
	"regexp"
	"strings"
)

func deleteHTML(s string) string {
	if strings.Contains(s, "<a") {
		re1, _ := regexp.Compile(`<[\S\s]+?>`)
		s = re1.ReplaceAllString(s, "")
		s = strings.Replace(s, "@", "+", -1)
	}
	return s
}

func Jiandan(messages chan<- tgbotapi.Chattable, comments <-chan types.Comment) {

	funcMap := template.FuncMap{"deleteHTML": deleteHTML}

	// todo 处理每一条帖子,然后发送
	for comment := range comments {

		// buffer用于接收模板输出
		var commentText bytes.Buffer

		commentTemplateText := `<a href="https://jandan.net/t/{{.Id}}">原帖链接</a>
{{.Author}}(楼主):{{.ContentText}}
OO:{{.OO}} XX:{{.XX}}
{{range .TuCao}}{{.Author}}:{{.Content|deleteHTML}}
OO:{{.OO}} XX:{{.XX}}
{{end}}`
		commentTemplate, _ := template.New("comment").Funcs(funcMap).Parse(commentTemplateText)

		err := commentTemplate.Execute(&commentText, comment)
		if err != nil {
			log.Println(err)
			continue
		}

		var medias []interface{}

		textAdded := false
		for _, pic := range comment.Pics {
			if textAdded {
				if pic[len(pic)-3:] != "gif" {
					medias = append(medias, tgbotapi.NewInputMediaPhoto(pic))
				} else {
					medias = append(medias, tgbotapi.NewInputMediaVideo(pic))
				}
			} else {
				if pic[len(pic)-3:] != "gif" {
					medias = append(medias, tgbotapi.InputMediaPhoto{
						Type:      "photo",
						Media:     pic,
						Caption:   commentText.String(),
						ParseMode: tgbotapi.ModeHTML,
					})
				} else {
					medias = append(medias, tgbotapi.InputMediaVideo{
						Type:      "video",
						Media:     pic,
						Caption:   commentText.String(),
						ParseMode: tgbotapi.ModeHTML,
					})
				}
				textAdded = true
			}
		}

		newMsg := tgbotapi.MediaGroupConfig{
			BaseChat: tgbotapi.BaseChat{
				ChannelUsername: viper.GetString("ChannelUsername"),
			},
			InputMedia: medias,
		}
		commentText.Reset()
		messages <- newMsg

	}
}
