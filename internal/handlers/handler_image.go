package handlers

import (
	"fmt"
	"strings"

	"github.com/xbmlz/chatgpt-plus-dingtalk/internal/config"
	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/dingbot"
	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/logger"
	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/replicate"
)

func HandlerImage(ding *dingbot.DingBot, msg dingbot.DingBotReceiveMessage) {
	image := replicate.New(replicate.Replicate{
		BaseUrl:  config.Instance.ReplicateBaseUrl,
		ApiToken: config.Instance.ReplicateApiToken,
	})

	url, err := image.Generate(replicate.ImageGenerateRequest{
		Version: config.Instance.ReplicateModelVersion,
		Input: replicate.ImageGenerateRequestInput{
			Prompt: strings.ReplaceAll(msg.Text.Content, "图片", ""),
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
}
