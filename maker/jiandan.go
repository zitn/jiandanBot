package maker

import (
	"bytes"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"myTeleBot/types"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

func Jiandan(messages chan tgbotapi.Chattable, comments chan types.Comment) {
	// todo 处理每一条帖子,然后发送
	for comment := range comments {

		// 在buffer中组合消息
		var commentText bytes.Buffer
		// 第一行 原贴链接
		//commentText.WriteString("[原贴链接](https://jandan.net/t/")
		//commentText.WriteString(comment.Id)
		//commentText.WriteString(")\n")

		// 楼主发言以及OO和XX
		commentText.WriteString(comment.Author)
		commentText.WriteString("(楼主):")
		if comment.ContentText == "" {
			commentText.WriteString("\n")
		} else {
			commentText.WriteString(comment.ContentText)
		}
		commentText.WriteString("OO:")
		commentText.WriteString(comment.OO)
		commentText.WriteString(" XX:")
		commentText.WriteString(comment.XX)

		// 若有吐槽,则加载
		if comment.SubCommentCount != "0" {
			for _, tucao := range comment.TuCao {
				commentText.WriteString("\n")
				commentText.WriteString(tucao.Author)
				commentText.WriteString(":")

				// 处理@人的情况
				if strings.Contains(tucao.Content, "<a href") {
					// 将html标签删除
					re1, _ := regexp.Compile(`<[\S\s]+?>`)
					tucao.Content = re1.ReplaceAllString(tucao.Content, "")
				}

				commentText.WriteString(tucao.Content)
				commentText.WriteString("\nOO:")
				commentText.WriteString(strconv.Itoa(tucao.OO))
				commentText.WriteString(" XX:")
				commentText.WriteString(strconv.Itoa(tucao.XX))
			}

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
						Type:    "photo",
						Media:   pic,
						Caption: commentText.String(),
					})
				} else {
					medias = append(medias, tgbotapi.InputMediaVideo{
						Type:    "video",
						Media:   pic,
						Caption: commentText.String(),
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
