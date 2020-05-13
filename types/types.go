package types

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type PageResult struct {
	Status        string    `json:"status"`
	CurrentPage   int       `json:"current_page"`
	TotalComments int       `json:"total_comments"`
	PageCount     int       `json:"page_count"`
	Count         int       `json:"count"`
	Comments      []Comment `json:"comments"`
}

type Comment struct {
	Id              string        `json:"comment_ID"`
	PostId          string        `json:"comment_post_ID"`
	Author          string        `json:"comment_author"`
	Date            string        `json:"comment_date"`
	DateGmt         string        `json:"comment_date_gmt"`
	Content         string        `json:"comment_content"`
	UserId          string        `json:"user_id"`
	OO              string        `json:"vote_positive"`
	XX              string        `json:"vote_negative"`
	SubCommentCount string        `json:"sub_comment_count"`
	ContentText     string        `json:"text_content"`
	Pics            []string      `json:"pics"`
	TuCao           []TuCadDetail `json:"-"`
}

type TuCadDetail struct {
	Id         int          `json:"id"`
	PostId     int          `json:"post_id"`
	Author     string       `json:"author"`
	AuthorType int          `json:"author_type"`
	Date       string       `json:"date"`
	AtComments []TuCaoAt    `json:"at_comments"`
	Content    string       `json:"content"`
	UserId     int          `json:"user_id"`
	XX         int          `json:"vote_negative"`
	OO         int          `json:"vote_positive"`
	Images     []TuCaoImage `json:"images"`
}

type TuCaoImage struct {
	Url       string `json:"url"`
	FullUrl   string `json:"full_url"`
	Host      string `json:"host"`
	ThumbSize string `json:"thumb_size"`
	Ext       string `json:"ext"`
	FileName  string `json:"file_name"`
}

type TuCaoAt struct {
	AtAuthor    string `json:"at_author"`
	AtCommentId string `json:"at_comment_id"`
}

type CommentMessage struct {
	HaveTucao      bool
	CommentMessage tgbotapi.MediaGroupConfig
	TucaoMessage   tgbotapi.MessageConfig
}

type TucaoUpdate struct {
	CommentId  string
	UpdateData tgbotapi.Update
}
