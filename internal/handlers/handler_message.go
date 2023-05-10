package handlers

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/xbmlz/chatgpt-plus-dingtalk/internal/config"
	"github.com/xbmlz/chatgpt-plus-dingtalk/internal/db"
	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/chatgpt"
	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/dingbot"
	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/logger"
)

func HandlerMessage(ding *dingbot.DingBot, msg dingbot.DingBotReceiveMessage) {
	qMessageID := uuid.NewString()
	db.Save(&db.Chat{
		DingTalkConversationID: msg.ConversationID,
		SenderID:               msg.SenderID,
		SenderNick:             msg.SenderNick,
		MessageID:              qMessageID,
		ConversationID:         "",
		ConversationMode:       db.QUESTION_MODE,
		ConversationType:       msg.ConversationType,
		Content:                msg.Text.Content,
	})

	var chatQuery db.Chat
	db.FindOne(map[string]interface{}{
		"ding_talk_conversation_id": msg.ConversationID,
		"conversation_type":         msg.ConversationType,
		"conversation_mode":         db.ANSWER_MODE,
	}, &chatQuery)

	var c chatgpt.CompletionRequest
	c.Action = "next"
	c.Messages = []chatgpt.CompletionRequestMessage{
		{
			ID:   qMessageID,
			Role: "user",
			Content: chatgpt.CompletionMessageContent{
				ContentType: "text",
				Parts:       []string{msg.Text.Content},
			},
		},
	}
	c.Model = config.Instance.ChatgptModel
	if chatQuery.ID > 0 {
		c.ConversationID = chatQuery.ConversationID
		c.ParentMessageID = chatQuery.MessageID
	} else {
		c.ConversationID = ""
		c.ParentMessageID = uuid.NewString()
	}
	// create completion
	chatgpt := chatgpt.New(chatgpt.ChatGPT{
		BaseUrl:     config.Instance.ChatgptBaseUrl,
		AccessToken: config.Instance.ChatgptAccessToken,
	})
	resp, err := chatgpt.CreateCompletion(c)
	if err != nil {
		logger.Error(err)
		errMsg := fmt.Sprintf("请求聊天机器人失败: %s", err.Error())
		ding.SendMessage(dingbot.MSG_TEXT, errMsg)
		return
	}
	respContent := resp.Message.Content.Parts[0]
	// send message
	err = ding.SendMessage(dingbot.MSG_MD, respContent)
	if err != nil {
		logger.Error(err)
	}
	// save answer message
	db.Save(&db.Chat{
		DingTalkConversationID: msg.ConversationID,
		SenderID:               msg.SenderID,
		SenderNick:             msg.SenderNick,
		MessageID:              resp.Message.ID,
		ConversationID:         resp.ConversationID,
		ConversationMode:       db.ANSWER_MODE,
		ConversationType:       msg.ConversationType,
		Content:                respContent,
	})
}
