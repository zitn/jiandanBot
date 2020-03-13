package maker

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"myTeleBot/channel"
	"myTeleBot/crawler"
	"myTeleBot/types"
	"strconv"
	"strings"
)

//var (
//	funcMap = template.FuncMap{"deleteHTML": deleteHTML}
//
//	// 楼主发言模板
//	commentTemplateText = `[原帖链接](https://jandan.net/t/{{.Id}})
//{{.Author}}:{{.ContentText}}
//OO:{{.OO}} XX:{{.XX}}`
//	commentTemplate, _ = template.New("comment").Funcs(funcMap).Parse(commentTemplateText)
//
//	// 吐槽模板
//	tucaoTemplateText = `{{range .}}{{.Author}}:{{.Content|deleteHTML}}
//OO:{{.OO}} XX:{{.XX}}
//{{end}}`
//	tucaoTemplate, _ = template.New("tucao").Funcs(funcMap).Parse(tucaoTemplateText)
//)

func Jiandan() {
	//  处理每一条帖子,然后发送
	for comment := range channel.CommentsChannel {

		var caption strings.Builder
		caption.WriteString("[原帖链接](https://jandan.net/t/")
		caption.WriteString(comment.Id)
		caption.WriteString(") By ")
		caption.WriteString(comment.Author)
		caption.WriteString("\n")
		caption.WriteString(comment.ContentText)
		caption.WriteString("\nOO:")
		caption.WriteString(comment.OO)
		caption.WriteString("   XX:")
		caption.WriteString(comment.XX)

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
						Caption:   caption.String(),
						ParseMode: tgbotapi.ModeMarkdown,
					})
				} else {
					medias = append(medias, tgbotapi.InputMediaVideo{
						Type:      "video",
						Media:     pic,
						Caption:   caption.String(),
						ParseMode: tgbotapi.ModeMarkdown,
					})
				}
				textAdded = true
			}
		}

		newComment := tgbotapi.MediaGroupConfig{
			BaseChat: tgbotapi.BaseChat{
				ChannelUsername: viper.GetString("ChannelUsername"),
			},
			InputMedia: medias,
		}

		// 吐槽
		tuCao := "========吐槽========"
		if comment.SubCommentCount != "0" {
			tuCao = generateTuCao(comment.TuCao)
		}

		// 更新按娘
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
			Text:                  tuCao,
			DisableWebPagePreview: true,
		}

		newMessage := types.CommentMessage{
			HaveTucao:      comment.SubCommentCount == "0",
			CommentMessage: newComment,
			TucaoMessage:   newTucao,
		}
		channel.CommentMessageChannel <- newMessage

	}
}

func generateTuCao(tuCaoDetail []types.TuCadDetail) string {
	var tuCaoBuilder strings.Builder
	for _, detail := range tuCaoDetail {
		// 如果吐槽中有图片,将图片链接添加进去
		for _, imageLink := range detail.Images {
			detail.Content = strings.Replace(detail.Content, `#img#`, " [图片]("+imageLink.FullUrl+") ", 1)
		}
		// 如果吐槽中at了别人,将其替换为 +username 的形式
		for _, atComment := range detail.AtComments {
			detail.Content = strings.Replace(detail.Content, `#at#`, " +*"+atComment.AtAuthor+"*", 1)
		}
		tuCaoBuilder.WriteString("*")
		tuCaoBuilder.WriteString(detail.Author)
		tuCaoBuilder.WriteString("*: ")
		tuCaoBuilder.WriteString(detail.Content)
		tuCaoBuilder.WriteString("\nOO:")
		tuCaoBuilder.WriteString(strconv.Itoa(detail.OO))
		tuCaoBuilder.WriteString("  XX:")
		tuCaoBuilder.WriteString(strconv.Itoa(detail.XX))
		tuCaoBuilder.WriteString("\n")
	}
	return tuCaoBuilder.String()
}

func UpdateTucao() {
	for req := range channel.RequireUpdateTucaoChannel {
		newTucaoDetails := crawler.GetTucao(req.CommentId)
		if len(newTucaoDetails) == 0 {
			continue
		}
		newTuCaoStr := generateTuCao(newTucaoDetails)
		if len(newTuCaoStr) < len(req.UpdateData.Message.Text)-20 {
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
			Text:                  newTuCaoStr,
			DisableWebPagePreview: true,
		}
		channel.NormalMessageChannel <- editedMsg
	}
}
