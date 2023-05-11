package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/dingbot"
	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/logger"
)

func RootHandler(ctx *gin.Context) {
	var msg dingbot.DingBotReceiveMessage
	err := ctx.Bind(&msg)
	if err != nil {
		return
	}
	ding := dingbot.New(msg)
	input := msg.Text.Content
	var ret string
	if strings.HasPrefix(input, "帮助") || input == "" {
		ret = HandlerHelp(msg)
	} else if strings.HasPrefix(input, "图片") {
		ret = HandlerImage(msg)
	} else if strings.HasPrefix(input, "流程图") {
		msg.Text.Content = "帮我使用mermaid-js 10.x中的graph TD语法设计一个 " + msg.Text.Content + " 流程图，并在最后做简要描述"
		ret = HandlerMermaid(msg)
	} else if strings.HasPrefix(input, "脑图") {
		msg.Text.Content = "帮我使用mermaid-js 10.x中的graph LR语法设计一个 " + msg.Text.Content + " 脑图，并在最后做简要描述"
		ret = HandlerMermaid(msg)
	} else {
		ret = HandlerMessage(msg)
	}
	err = ding.SendMessage(dingbot.MSG_MD, ret)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "ok"})
}
