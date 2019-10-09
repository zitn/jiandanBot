package crawler

import (
	"encoding/json"
	"errors"
	"github.com/spf13/viper"
	"io/ioutil"
	"myTeleBot/channel"
	"myTeleBot/types"
	"net/http"
	"time"
)

var (
	myClient = &http.Client{Timeout: 30 * time.Second}
)

func GetJiandan() {
	lastCommentTime := time.Now().Add(-time.Hour)
	var newTime time.Time
	tmpTime := lastCommentTime
	for {
		comments, err := getCommentList()
		if err != nil {
			channel.ErrorMessage <- err
			//log.Println(err)
			continue
		}
		for _, comment := range comments {
			newTime, err = time.Parse("2006-01-02 15:04:05", comment.Date)
			if err != nil {
				//log.Println(err)
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
					comment.TuCao, err = GetTucao(comment.Id)
					if err != nil {
						//log.Println(err)
						channel.ErrorMessage <- err
						continue
					}
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

func getCommentList() ([]types.Comment, error) {
	r, err := myClient.Get(viper.GetString("ApiAddress"))
	if err != nil {
		return nil, errors.New("time out")
	}
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, errors.New("read body error")
	}
	var commentList types.CommentList
	err = json.Unmarshal(body, &commentList)
	if err != nil {
		return nil, errors.New("unmarshal error")
	}
	err = r.Body.Close()
	if err != nil {
		return nil, errors.New("can't close client")
	}
	return commentList.Comments, nil
}

func GetTucao(commentID string) ([]types.TuCaoDetial, error) {
	r, err := myClient.Get("https://jandan.net/tucao/all/" + commentID)
	if err != nil {
		return nil, errors.New("read body error")
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.New("read body error")
	}
	err = r.Body.Close()
	if err != nil {
		return nil, errors.New("can't close client")
	}
	var tuCao types.TuCao
	err = json.Unmarshal(body, &tuCao)
	if err != nil {
		return nil, errors.New("unmarshal error")
	}
	return tuCao.Tucao, nil
}
