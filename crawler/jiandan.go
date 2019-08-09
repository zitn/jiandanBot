package crawler

import (
	"encoding/json"
	"errors"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"myTeleBot/channel"
	"myTeleBot/types"
	"net/http"
	"time"
)

var (
	myClient        = &http.Client{Timeout: 30 * time.Second}
	lastCommentTime = time.Now().Add(30 * time.Minute) // todo 30分钟为测试之用
)

func init() {
	go getJiandan()
}

func getJiandan() {
	for {
		comments, err := getCommentList()
		if err != nil {
			log.Println(err)
			continue
		}
		for _, comment := range comments {
			// 如果为新帖子,获取吐槽,将数据发送给maker进行处理
			newTime, err := time.Parse("2006-01-02 15:04:05", comment.Date)
			if err != nil {
				log.Println(err)
				continue
			}
			if lastCommentTime.Before(newTime) {
				// 如果新帖子的吐槽数不为0,则获取吐槽
				if comment.SubCommentCount != "0" {
					comment.TuCao, err = GetTucao(comment.Id)
					if err != nil {
						log.Println(err)
						continue
					}
				}
				channel.CommentsChannel <- comment
			} else {
				// 如果id重复,则停止发送剩余数据
				lastCommentTime = newTime
				break
			}
		}
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
	r, err := myClient.Get("https://i.jandan.net/tucao/" + commentID)
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
