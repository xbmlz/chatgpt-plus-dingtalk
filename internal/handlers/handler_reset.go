package handlers

import (
	"github.com/xbmlz/chatgpt-plus-dingtalk/internal/config"
	"github.com/xbmlz/chatgpt-plus-dingtalk/internal/db"
	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/chatgpt"
	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/dingbot"
)

func HandlerReset(msg dingbot.DingBotReceiveMessage) (retMsg string) {
	chatgpt := chatgpt.New(chatgpt.ChatGPT{
		BaseUrl:     config.Instance.ChatgptBaseUrl,
		AccessToken: config.Instance.ChatgptAccessToken,
	})
	err := chatgpt.DeleteConversation(msg.ConversationID)
	if err != nil {
		err = db.DeleteByDingTalkConversationID(msg.ConversationID)
		if err != nil {
			retMsg = "❌ 重置会话失败, 请稍后再试"
		}
	}
	retMsg = "♻️ 重置会话成功"
	return
}
