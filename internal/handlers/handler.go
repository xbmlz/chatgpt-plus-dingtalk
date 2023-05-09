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
	"github.com/xbmlz/chatgpt-dingtalk/pkg/replicate"
)

func RootHandler(ctx *gin.Context) {
	var msg dingbot.DingBotReceiveMessage
	err := ctx.Bind(&msg)
	if err != nil {
		return
	}
	ding := dingbot.New(msg)
	// TODO
	input := msg.Text.Content
	if strings.HasPrefix(input, "帮助") {
		SendHelp(ding)
		return
	}
	if strings.HasPrefix(input, "#图片") {
		input = strings.ReplaceAll(input, "#图片", "")
		image := replicate.New(replicate.Replicate{
			BaseUrl:  "https://api.replicate.com",
			ApiToken: "r8_EFqWf2Io13JHek548wcRb3I8Z34KDM51ykCia",
		})

		url, err := image.Generate(replicate.ImageGenerateRequest{
			Version: "db21e45d3f7023abc2a46ee38a23973f6dce16bb082a930b0c49861f96d1e5bf",
			Input: replicate.ImageGenerateRequestInput{
				Prompt: input,
			},
		})
		if err != nil {
			logger.Error(err)
			errMsg := fmt.Sprintf("请求聊天机器人失败: %s", err.Error())
			ding.SendMessage(dingbot.MSG_TEXT, errMsg)
			return
		}
		imgMd := fmt.Sprintf("![image](%s)", url)
		ding.SendMessage(dingbot.MSG_MD, imgMd)
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

func SendHelp(ding *dingbot.DingBot) {
	content := `
	### 🤖 需要帮助吗？

	我是卫博士，一款基于ChatGPT技术的智能聊天机器人！
	
	回复 **图片 + 描述** 或 **/img + 描述** 生成图片。
	回复 **帮助** 或 **help** 获取帮助信息。

	`
	ding.SendMessage(dingbot.MSG_MD, content)
}
