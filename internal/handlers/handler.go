package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/dingbot"
	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/logger"
)

type ActionType string

func RootHandler(ctx *gin.Context) {
	var msg dingbot.DingBotReceiveMessage
	err := ctx.Bind(&msg)
	if err != nil {
		return
	}
	msg.MsgID = uuid.NewString()
	ding := dingbot.New(msg)
	input := msg.Text.Content
	var retMsg string
	if strings.HasPrefix(input, "帮助") || input == "" {
		retMsg = HandlerHelp(msg)
	} else if strings.HasPrefix(input, "图片") {
		retMsg = HandlerImage(msg)
	} else if strings.HasPrefix(input, "流程图") {
		retMsg = HandlerFlowchart(msg)
	} else if strings.HasPrefix(input, "脑图") {
		retMsg = HandlerMindmap(msg)
	} else if strings.HasPrefix(input, "重置") {
		retMsg = HandlerReset(msg)
	} else {
		retMsg = HandlerMessage(msg)
	}
	err = ding.SendMessage(dingbot.MSG_MD, retMsg)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "ok"})
}
