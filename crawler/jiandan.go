package crawler

import (
	"github.com/go-resty/resty/v2"
	"github.com/json-iterator/go"
	"github.com/spf13/viper"
	"myTeleBot/channel"
	"myTeleBot/types"
	"net/http"
	"time"
)

var (
	request = resty.New()
	json    = jsoniter.ConfigFastest
)

func GetJianDan() {
	lastCommentTime := time.Now().Add(-time.Minute * 5)
	var newTime time.Time
	tmpTime := lastCommentTime
	for {
		comments, err := getNewComments()
		if err != nil {
			channel.ErrorMessage <- err
			continue
		}
		for _, comment := range comments {
			newTime, err = time.Parse("2006-01-02 15:04:05", comment.Date)
			if err != nil {
				channel.ErrorMessage <- err
				continue
			}
			if tmpTime.Before(newTime) {
				tmpTime = newTime
			}
			// 如果记录的最晚帖子时间在新帖子时间之前
			if lastCommentTime.Before(newTime) {
				// 如果新帖子的吐槽数不为0,则获取吐槽
				if comment.SubCommentCount != "0" {
					// todo 获取吐槽失败的错误在函数内部处理
					comment.TuCao = GetTucao(comment.Id)
				}
				channel.CommentsChannel <- comment
			} else {
				// 则停止发送剩余数据
				break
			}
		}

		lastCommentTime = tmpTime
		// 20分钟get一次数据
		time.Sleep(20 * time.Minute)

	}
}

func getNewComments() ([]types.Comment, error) {
	response, err := request.R().Get(viper.GetString("ApiAddress"))
	if err != nil {
		return nil, err
	}
	var comments []types.Comment
	json.Get(response.Body(), "comments").ToVal(&comments)
	return comments, nil
}

func GetTucao(commentID string) []types.TuCadDetail {
	response, err := request.R().Get("https://api.jandan.net/api/v1/tucao/list/" + commentID)
	if err != nil {
		return nil
	}
	if response.StatusCode() != http.StatusOK {
		return nil
	}
	if response.Body() == nil {
		return nil
	}
	var TucaoDetails []types.TuCadDetail
	json.Get(response.Body(), "data").ToVal(&TucaoDetails)
	return TucaoDetails
}
