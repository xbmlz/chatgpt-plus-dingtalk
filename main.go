package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xbmlz/chatgpt-dingtalk/internal/config"
	"github.com/xbmlz/chatgpt-dingtalk/internal/db"
	"github.com/xbmlz/chatgpt-dingtalk/pkg/chatgpt"
	"github.com/xbmlz/chatgpt-dingtalk/pkg/dingbot"
	"github.com/xbmlz/chatgpt-dingtalk/pkg/logger"
)

// dingtalk bot

// chatgpt

func main() {
	config.Init()
	logger.Init(config.Instance.LogLevel)
	db.Init()
	r := gin.Default()
	r.POST("/", func(ctx *gin.Context) {
		var msg dingbot.DingBotReceiveMessage
		err := ctx.Bind(&msg)
		if err != nil {
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
		c.Messages = append(c.Messages, chatgpt.CompletionRequestMessage{
			ID:   qMessageID,
			Role: "user",
			Content: chatgpt.CompletionMessageContent{
				ContentType: "text",
				Parts:       []string{msg.Text.Content},
			},
		})
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
		}
		respContent := resp.Message.Content.Parts[0]
		// send message
		dingbot := dingbot.NewDingBot(dingbot.DingBot{
			SessionWebhook: msg.SessionWebhook,
		})
		err = dingbot.SendMarkdownMessage(msg.ConversationType, "markdown", respContent, msg.SenderStaffID, msg.SenderNick)
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
	})
	port := ":" + config.Instance.ServerPort
	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		logger.Info("🚀 Listening and serving HTTP on", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Listen: %s\n", err)
		}
	}()
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info("Shutdown Server ...")

	// 5秒后强制退出
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server Shutdown: %s", err)
	}
	logger.Info("Server exiting!")
}
