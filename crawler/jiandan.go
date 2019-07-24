package crawler

import (
	"encoding/json"
	"errors"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"myTeleBot/types"
	"net/http"
	"time"
)

var (
	myClient    = &http.Client{Timeout: 30 * time.Second}
	lastComment string
)

func GetJiandan(commentsChan chan<- types.Comment) {
	for {
		comments, err := getCommentList(viper.GetString("ApiAddress"))
		if err != nil {
			log.Println(err)
			continue
		}
		for _, comment := range comments {

			if comment.Id != lastComment {
				// 如果为新帖子,获取吐槽,将数据发送给maker进行处理
				if comment.SubCommentCount != "0" {
					comment.TuCao, err = getTucao("https://i.jandan.net/tucao/" + comment.Id)
					if err != nil {
						log.Println(err)
						continue
					}
				}
				commentsChan <- comment
			} else {
				// 如果id重复,则停止发送剩余数据
				break
			}
		}

		lastComment = comments[0].Id
		// 20分钟get一次数据
		time.Sleep(20 * time.Minute)

	}
}

// todo 重构get函数

func getCommentList(url string) ([]types.Comment, error) {
	r, err := myClient.Get(url)
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

func getTucao(url string) ([]types.TuCaoDetial, error) {
	r, err := myClient.Get(url)
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
