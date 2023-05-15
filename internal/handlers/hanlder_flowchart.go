package handlers

import (
	"fmt"

	"github.com/xbmlz/chatgpt-plus-dingtalk/internal/config"
	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/dingbot"
	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/utils"
)

func HandlerFlowchart(msg dingbot.DingBotReceiveMessage) (retMsg string) {
	msg.Text.Content = "请你用 mermaid 语法帮我完成一个关于 " + msg.Text.Content + " 的流程图, 并在节点上添加中文描述信息"
	retID, retMsg := AskChatGPT(msg)
	between := utils.ExtractStringBetween(retMsg, "```mermaid", "```")
	if between != "" {
		link := fmt.Sprintf(config.Instance.ServerUrl+"/blob?id=%s&type=flowchart", retID)
		retMsg = fmt.Sprintf("\n\n[✔️ 生成流程图成功, 点击查看](%s)", link)
	} else {
		retMsg = "❌ 生成流程图失败, 请稍后再试"
	}
	return
}
