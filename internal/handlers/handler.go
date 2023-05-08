package handlers

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xbmlz/chatgpt-dingtalk/internal/config"
	"github.com/xbmlz/chatgpt-dingtalk/internal/db"
	"github.com/xbmlz/chatgpt-dingtalk/pkg/chatgpt"
	"github.com/xbmlz/chatgpt-dingtalk/pkg/dingbot"
	"github.com/xbmlz/chatgpt-dingtalk/pkg/logger"
)

func RootHandler(ctx *gin.Context) {
	var msg dingbot.DingBotReceiveMessage
	err := ctx.Bind(&msg)
	if err != nil {
		return
	}
	ding := dingbot.NewDingBot(msg)
	// TODO
	if strings.HasPrefix(msg.Text.Content, "帮助") {
		SendHelp(ding)
		return
	}

	// save question message
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
			ID:   uuid.NewString(),
			Role: "system",
			Content: chatgpt.CompletionMessageContent{
				ContentType: "text",
				Parts:       []string{"你是 ChatGPT，一个由 OpenAI 训练的大型语言模型。请仔细遵循用户的指示。使用 Markdown 格式进行回应。"},
			},
		},
		{
			ID:   qMessageID,
			Role: "user",
			Content: chatgpt.CompletionMessageContent{
				ContentType: "text",
				Parts:       []string{msg.Text.Content},
			},
		},
	}
	c.Model = config.Instance.Model
	if chatQuery.ID > 0 {
		c.ConversationID = chatQuery.ConversationID
		c.ParentMessageID = chatQuery.MessageID
	} else {
		c.ConversationID = ""
		c.ParentMessageID = uuid.NewString()
	}
	// create completion
	chatgpt := chatgpt.NewChatGPT(chatgpt.ChatGPT{
		BaseUrl:     config.Instance.ApiUrl,
		AccessToken: config.Instance.AccessToken,
	})
	resp, err := chatgpt.CreateCompletion(c)
	if err != nil {
		logger.Error(err)
		errMsg := fmt.Sprintf("请求聊天机器人失败: %s", err.Error())
		ding.SendMessage(dingbot.MSG_TEXT, errMsg)
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

func SendHelp(ding *dingbot.DingBot) {
	content := `
	### 🤖 需要帮助吗？

	我是卫博士，一款基于ChatGPT技术的智能聊天机器人！
	
	回复 **图片 + 描述** 或 **/img + 描述** 生成图片。
	回复 **帮助** 或 **help** 获取帮助信息。

	`
	ding.SendMessage(dingbot.MSG_MD, content)
}
