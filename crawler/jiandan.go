package crawler

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"jiandanBot/channel"
	"jiandanBot/types"
)

var (
	request = resty.New()
	json    = jsoniter.ConfigFastest
)

func init() {
	header := map[string]string{
		"accept":                    `text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9`,
		"accept-language":           `zh-CN,zh-TW;q=0.9,zh;q=0.8,en-US;q=0.7,en;q=0.6`,
		"sec-ch-ua":                 ` " Not;A Brand";v="99", "Google Chrome";v="91", "Chromium";v="91"`,
		"sec-ch-ua-mobile":          "?0",
		"sec-fetch-dest":            "document",
		"sec-fetch-mode":            "navigate",
		"sec-fetch-site":            "none",
		"sec-fetch-user":            "?1",
		"upgrade-insecure-requests": "1",
	}
	request.SetHeaders(header)
}

func GetJianDan() {
	log1 := logrus.WithField("func", "GetJianDan")
	var lastCommentTime time.Time
	configFileTime := viper.GetString("lastCommentTime")
	if configFileTime == "" {
		lastCommentTime = time.Now().Add(-time.Minute * 5)
	} else {
		var err error
		lastCommentTime, err = time.Parse("2006-01-02 15:04:05", configFileTime)
		if err != nil {
			log1.Error(err)
			lastCommentTime = time.Now().Add(-time.Minute * 5)
		}
	}
	var newTime time.Time
	tmpTime := lastCommentTime
	for {
		comments, err := getNewComments()
		if err != nil {
			log1.WithField("err in ", "getNewComments").Error(err)
			continue
		}
		for _, comment := range comments {
			newTime, err = time.Parse("2006-01-02 15:04:05", comment.Date)
			if err != nil {
				log1.WithField("err in ", "time.Parse").Error(err)
				continue
			}
			if tmpTime.Before(newTime) {
				tmpTime = newTime
			}
			// 如果记录的最晚帖子时间在新帖子时间之前
			if lastCommentTime.Before(newTime) {
				// 如果新帖子的吐槽数不为0,则获取吐槽
				if comment.SubCommentCount != "0" {
					comment.TuCao = GetTucao(comment.Id)
				}
				channel.CommentsChannel <- comment
			} else {
				// 处理sein的置顶帖
				if comment.Author == "sein" {
					continue
				}
				// 则停止发送剩余数据
				break
			}
		}

		lastCommentTime = tmpTime
		viper.Set("lastCommentTime", lastCommentTime.Format("2006-01-02 15:04:05"))
		// 20分钟get一次数据
		time.Sleep(15 * time.Minute)
	}
}

func getNewComments() ([]types.Comment, error) {
	response, err := request.R().Get(viper.GetString("ApiAddress"))
	if err != nil {
		return nil, err
	}
	if response.StatusCode() != 200 {
		return nil, errors.New("response error")
	}
	var pageResult types.PageResult
	if response.Body() == nil {
		return nil, errors.New("response is empty")
	}
	err = json.Unmarshal(response.Body(), &pageResult)
	if err != nil {
		return nil, err
	}
	return pageResult.Comments, nil
}

func GetTucao(commentID string) []types.TuCaoDetail {
	log1 := logrus.WithField("func", "GetTucao")
	response, err := request.R().Get("https://api.jandan.net/api/v1/tucao/list/" + commentID)
	if err != nil {
		log1.WithField("err in", "request").WithField("commentID", commentID).Error(err)
		return nil
	}
	if response.StatusCode() != http.StatusOK {
		log1.WithField("err in", "response.StatusCode").WithField("commentID", commentID).Error("response.StatusCode is", response.StatusCode())
		return nil
	}
	if response.Body() == nil {
		log1.WithField("err in", "response.Body").WithField("commentID", commentID).Error("response body is nil")
		return nil
	}
	var TucaoDetails []types.TuCaoDetail
	json.Get(response.Body(), "data", "list").ToVal(&TucaoDetails)
	return TucaoDetails
}
