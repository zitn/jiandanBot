// 全局 channel
package channel

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"myTeleBot/types"
)

var (
	// 普通消息
	NormalMessageChannel chan tgbotapi.Chattable
	// 煎蛋帖子
	CommentMessageChannel chan types.CommentMessage
	// 煎蛋爬虫爬取的原始数据
	CommentsChannel chan types.Comment
	// 煎蛋更新需求channel
	RequireUpdateTucaoChannel chan types.TucaoUpdate
	// 错误日志发送至管理员channel
	ErrorMessage chan error
)

// 初始化所有 channel
func init() {
	NormalMessageChannel = make(chan tgbotapi.Chattable, 100)
	CommentMessageChannel = make(chan types.CommentMessage, 40)
	CommentsChannel = make(chan types.Comment, 30)
	RequireUpdateTucaoChannel = make(chan types.TucaoUpdate, 30)
	ErrorMessage = make(chan error, 100)
}
