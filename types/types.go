package types

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type CommentList struct {
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
	TuCao           []TuCaoDetial `json:"-"`
}

type TuCao struct {
	Code        int           `json:"code"`
	HasNextPage bool          `json:"has_next_page"`
	HotTucao    []TuCaoDetial `json:"hot_tucao"`
	Tucao       []TuCaoDetial `json:"tucao"`
}

type TuCaoDetial struct {
	Id           int    `json:"comment_id"`
	Author       string `json:"comment_author"`
	Content      string `json:"comment_content"`
	Date         string `json:"comment_date"`
	DateInt      int    `json:"comment_date_int"`
	Parent       int    `json:"comment_parent"`
	PostId       int    `json:"comment_post_id"`
	ReplyId      int    `json:"comment_reply_id"`
	IsJandanUser int    `json:"is_jandan_user"`
	IsTipUser    int    `json:"is_tip_user"`
	XX           int    `json:"vote_negative"`
	OO           int    `json:"vote_positive"`
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
