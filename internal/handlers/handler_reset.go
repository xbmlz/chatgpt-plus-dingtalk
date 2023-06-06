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
	retMsg = "❌ 重置会话失败, 请稍后再试"
	gptConversationID, err := db.FindGptConversationId(msg.ConversationID)
	if err != nil {
		return
	}
	if gptConversationID != "" {
		err = chatgpt.DeleteConversation(gptConversationID)
		if err != nil {
			return
		}
	}
	err = db.DeleteByDingTalkConversationID(msg.ConversationID)
	if err != nil {
		return
	}

	retMsg = "♻️ 重置会话成功"
	return
}
