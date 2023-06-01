package task

import (
	"github.com/xbmlz/chatgpt-plus-dingtalk/internal/db"
	"github.com/xbmlz/chatgpt-plus-dingtalk/internal/handlers"
	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/dingbot"
	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/logger"
	"time"
)

func CleanAllSessionInterval() {
	logger.Info("会话清除定时任务启动", time.Now())

	conversationID, err := db.FindAllConversationID()
	if err != nil {
		logger.Error("查询会话ID失败：", err)
		return
	}
	for _, v := range conversationID {
		msg := dingbot.DingBotReceiveMessage{
			ConversationID: v,
		}
		handlers.HandlerReset(msg)
	}

	logger.Info("会话清除定时任务完成", time.Now())
}
