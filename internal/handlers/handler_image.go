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
	msg.Text.Content = "你帮我生成一个关于 " + msg.Text.Content + " 的图片, 用Unsplash API表示，并遵循以下的格式：https://source.unsplash.com/1600x900/?< PUT YOUR QUERY HERE >，只回复我链接，不需要其他表述"
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
	// _, retMsg = AskChatGPT(msg)
	// retMsg = fmt.Sprintf("![image](%s)", retMsg)
	return
}
