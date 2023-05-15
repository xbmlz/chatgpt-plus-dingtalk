package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xbmlz/chatgpt-plus-dingtalk/internal/db"
	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/utils"
)

func BlobHandler(ctx *gin.Context) {
	id := ctx.Query("id")
	t := ctx.Query("type")
	var chat db.Chat
	db.FindOne(map[string]interface{}{
		"message_id": id,
	}, &chat)

	if chat.ID > 0 {
		if t == "flowchart" {
			chat.Content = utils.ExtractStringBetween(chat.Content, "```mermaid", "```")
			chat.Content = strings.ReplaceAll(chat.Content, "LR", "TD")
		} else if t == "mindmap" {
			chat.Content = utils.ExtractStringBetween(chat.Content, "```markmap", "```")
		}
	}

	ctx.HTML(http.StatusOK, "blob.html", gin.H{
		"type": t,
		"chat": chat,
	})
}
