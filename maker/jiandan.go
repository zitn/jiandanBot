package maker

import (
	"bytes"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"html/template"
	"myTeleBot/channel"
	"myTeleBot/crawler"
	"myTeleBot/types"
	"regexp"
	"strings"
)

var (
	funcMap = template.FuncMap{"deleteHTML": deleteHTML}

	// 楼主发言模板
	commentTemplateText = `[原帖链接](https://jandan.net/t/{{.Id}})
{{.Author}}:{{.ContentText}}
OO:{{.OO}} XX:{{.XX}}`
	commentTemplate, _ = template.New("comment").Funcs(funcMap).Parse(commentTemplateText)

	// 吐槽模板
	tucaoTemplateText = `{{range .}}{{.Author}}:{{.Content|deleteHTML}}
OO:{{.OO}} XX:{{.XX}}
{{end}}`
	tucaoTemplate, _ = template.New("tucao").Funcs(funcMap).Parse(tucaoTemplateText)
)

func Jiandan() {
	//  处理每一条帖子,然后发送
	for comment := range channel.CommentsChannel {
		var commentBuff bytes.Buffer

		err := commentTemplate.Execute(&commentBuff, comment)
		if err != nil {
			//log.Println(err)
			channel.ErrorMessage <- err
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
						Caption:   commentBuff.String(),
						ParseMode: tgbotapi.ModeMarkdown,
					})
				} else {
					medias = append(medias, tgbotapi.InputMediaVideo{
						Type:      "video",
						Media:     pic,
						Caption:   commentBuff.String(),
						ParseMode: tgbotapi.ModeMarkdown,
					})
				}
				textAdded = true
			}
		}
		commentBuff.Reset()

		newComment := tgbotapi.MediaGroupConfig{
			BaseChat: tgbotapi.BaseChat{
				ChannelUsername: viper.GetString("ChannelUsername"),
			},
			InputMedia: medias,
		}

		// 吐槽
		var tucaoBuff bytes.Buffer
		tucaoBuff.WriteString("=======吐槽=======\n")
		if comment.SubCommentCount != "0" {

			err = tucaoTemplate.Execute(&tucaoBuff, comment.TuCao)
			if err != nil {
				//log.Println(err)
				channel.ErrorMessage <- err
				continue
			}
		}

		numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("更新吐槽", "updateTucao "+comment.Id),
			),
		)

		newTucao := tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChannelUsername:     viper.GetString("ChannelUsername"),
				DisableNotification: true,
				ReplyMarkup:         numericKeyboard,
			},
			ParseMode:             tgbotapi.ModeMarkdown,
			Text:                  tucaoBuff.String(),
			DisableWebPagePreview: true,
		}
		tucaoBuff.Reset()

		newMessage := types.CommentMessage{
			HaveTucao:      comment.SubCommentCount == "0",
			CommentMessage: newComment,
			TucaoMessage:   newTucao,
		}
		channel.CommentMessageChannel <- newMessage

	}
}

func UpdateTucao() {
	for req := range channel.RequireUpdateTucaoChannel {
		newTucao, err := crawler.GetTucao(req.CommentId)
		if err != nil {
			//log.Println(err)
			channel.ErrorMessage <- err
			continue
		}
		if len(newTucao) == 0 {
			//处理无吐槽的情况
			continue
		}
		var tucaoBuff bytes.Buffer
		tucaoBuff.WriteString("=======吐槽=======\n")
		err = tucaoTemplate.Execute(&tucaoBuff, newTucao)
		if err != nil {
			//log.Println(err)
			channel.ErrorMessage <- err
			continue
		}
		numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("更新吐槽", "updateTucao "+req.CommentId),
			),
		)

		editedMsg := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				//ChatID:          req.UpdateData.CallbackQuery.Message.Chat.ID,
				ChannelUsername: viper.GetString("ChannelUsername"),
				MessageID:       req.UpdateData.CallbackQuery.Message.MessageID,
				ReplyMarkup:     &numericKeyboard,
			},
			Text:                  tucaoBuff.String(),
			DisableWebPagePreview: true,
		}
		tucaoBuff.Reset()
		channel.NormalMessageChannel <- editedMsg
	}
}

// 删除@人标签的过滤器
func deleteHTML(s string) string {
	if strings.Contains(s, "<a") {
		re1, _ := regexp.Compile(`<[\S\s]+?>`)
		s = re1.ReplaceAllString(s, "")
		s = strings.Replace(s, "@", "&#43;", -1)
	}
	return s
}
