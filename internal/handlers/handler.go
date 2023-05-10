package handlers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/dingbot"
)

func RootHandler(ctx *gin.Context) {
	var msg dingbot.DingBotReceiveMessage
	err := ctx.Bind(&msg)
	if err != nil {
		return
	}
	ding := dingbot.New(msg)
	if strings.HasPrefix(msg.Text.Content, "帮助") {
		HandlerHelp(ding)
		return
	} else if strings.HasPrefix(msg.Text.Content, "图片") {
		HandlerImage(ding, msg)
		return
	} else {
		HandlerMessage(ding, msg)
		return
	}
}
