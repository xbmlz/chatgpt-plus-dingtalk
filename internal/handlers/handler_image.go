package handlers

import (
	"fmt"
	"strings"

	"github.com/xbmlz/chatgpt-plus-dingtalk/internal/config"
	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/dingbot"
	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/logger"
	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/replicate"
)

func HandlerImage(msg dingbot.DingBotReceiveMessage) (retMsg string) {
	image := replicate.New(replicate.Replicate{
		BaseUrl:  config.Instance.ReplicateBaseUrl,
		ApiToken: config.Instance.ReplicateApiToken,
	})
	prompt := strings.ReplaceAll(msg.Text.Content, "图片", "")
	url, err := image.Generate(replicate.ImageGenerateRequest{
		Version: config.Instance.ReplicateModelVersion,
		Input: replicate.ImageGenerateRequestInput{
			Prompt: prompt,
		},
	})
	if err != nil {
		logger.Error(err)
		retMsg = fmt.Sprintf("🚨 replicate 请求失败，请联系管理员: %s", err.Error())
		// ding.SendMessage(dingbot.MSG_TEXT, errMsg)
		return
	}
	retMsg = fmt.Sprintf("![image](%s)", url)
	// ding.SendMessage(dingbot.MSG_MD, imgMd)
	return
}
