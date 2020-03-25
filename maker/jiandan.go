package maker

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"jiandanBot/channel"
	"jiandanBot/crawler"
	"jiandanBot/types"
	"strconv"
	"strings"
)

func Jiandan() {
	//  处理每一条帖子,然后发送
	var caption strings.Builder
	for comment := range channel.CommentsChannel {

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
		tuCao := "========暂无吐槽========"
		if comment.SubCommentCount != "0" {
			tmp := generateTuCao(comment.TuCao)
			if tmp == "" {
				logrus.Error("generateTuCao", comment.TuCao)
				tuCao = tmp
			}
		}

		// 更新吐槽按钮
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
		// todo remove dev code
		caption.Reset()
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

func UpdateTucao(req tgbotapi.Update, commentID string) string {
	if req.CallbackQuery.Message == nil {
		fmt.Println(req)
		return "本条吐槽无法被更新"
	}
	newTucaoDetails := crawler.GetTucao(commentID)
	if len(newTucaoDetails) == 0 {
		return "信号消失在了火星"
	}
	newTuCaoStr := generateTuCao(newTucaoDetails)
	if len(newTuCaoStr) < len(req.CallbackQuery.Message.Text)-10 {
		return "没有更多吐槽"
	}

	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("更新吐槽", "updateTucao "+commentID),
		),
	)

	editedMsg := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChannelUsername: viper.GetString("ChannelUsername"),
			MessageID:       req.CallbackQuery.Message.MessageID,
			ReplyMarkup:     &numericKeyboard,
		},
		Text:                  newTuCaoStr,
		DisableWebPagePreview: true,
		ParseMode:             tgbotapi.ModeMarkdown,
	}
	channel.NormalMessageChannel <- editedMsg
	return "已更新"
}
