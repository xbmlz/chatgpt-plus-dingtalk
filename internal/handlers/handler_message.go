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

func HandlerMessage(msg dingbot.DingBotReceiveMessage) (retMsg string) {
	retID, retMsg := AskChatGPT(msg)
	link := fmt.Sprintf(config.Instance.ServerUrl+"/blob?id=%s&type=markdown", retID)
	retMsg += fmt.Sprintf("\n\n[在浏览器中查看此消息](%s)", link)
	return
}

func AskChatGPT(msg dingbot.DingBotReceiveMessage) (retID, retMsg string) {
	db.Save(&db.Chat{
		DingTalkConversationID: msg.ConversationID,
		SenderID:               msg.SenderID,
		SenderNick:             msg.SenderNick,
		MessageID:              msg.MsgID,
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
			ID:   msg.MsgID,
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
		retMsg = fmt.Sprintf("chatgpt请求失败，请联系管理员: %s", err.Error())
		return
	}
	retMsg = resp.Message.Content.Parts[0]
	retID = resp.Message.ID
	db.Save(&db.Chat{
		DingTalkConversationID: msg.ConversationID,
		SenderID:               msg.SenderID,
		SenderNick:             msg.SenderNick,
		MessageID:              resp.Message.ID,
		ConversationID:         resp.ConversationID,
		ConversationMode:       db.ANSWER_MODE,
		ConversationType:       msg.ConversationType,
		Content:                retMsg,
	})
	return
}
