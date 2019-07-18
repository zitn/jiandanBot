package crawler

import (
	"encoding/json"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"myTeleBot/types"
	"net/http"
	"time"
)

var (
	myClient    = &http.Client{Timeout: 10 * time.Second}
	lastComment string
)

func GetJiandan(commentsChan chan types.Comment) {
	for {
		comments := getCommentList(viper.GetString("ApiAddress"))
		for _, comment := range comments {

			if comment.Id != lastComment {
				// 如果为新帖子,获取吐槽,将数据发送给maker进行处理
				if comment.SubCommentCount != "0" {
					comment.TuCao = getTucao("https://i.jandan.net/tucao/" + comment.Id)
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

		// todo 测试代码 5秒抓取一次数据
		//time.Sleep(5 * time.Second)
	}

}

func getCommentList(url string) []types.Comment {
	r, err := myClient.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Panic(err)
	}

	var commentList types.CommentList
	err = json.Unmarshal(body, &commentList)
	if err != nil {
		log.Panic(err)
	}
	return commentList.Comments
}

func getTucao(url string) []types.TuCaoDetial {
	r, err := myClient.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Panic(err)
	}

	var tuCao types.TuCao
	err = json.Unmarshal(body, &tuCao)
	if err != nil {
		log.Panic(err)
	}

	return tuCao.Tucao
}
