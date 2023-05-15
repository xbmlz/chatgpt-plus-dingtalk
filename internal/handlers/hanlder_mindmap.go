package handlers

import (
	"fmt"

	"github.com/xbmlz/chatgpt-plus-dingtalk/internal/config"
	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/dingbot"
	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/utils"
)

func HandlerMindmap(msg dingbot.DingBotReceiveMessage) (retMsg string) {
	msg.Text.Content = "请你用 markmap 语法帮我完成一个关于 " + msg.Text.Content + " 的思维导图, 并在节点上添加中文描述信息"
	retID, retMsg := AskChatGPT(msg)
	between := utils.ExtractStringBetween(retMsg, "```markmap", "```")
	if between != "" {
		link := fmt.Sprintf(config.Instance.ServerUrl+"/blob?id=%s&type=mindmap", retID)
		retMsg = fmt.Sprintf("\n\n[✔️ 生成脑图成功, 点击查看](%s)", link)
	} else {
		retMsg = "❌ 生成脑图失败, 请稍后再试"
	}

	// chatMsg := HandlerMessage(msg)
	// between := utils.ExtractStringBetween("```mermaid", "```", chatMsg)
	// // between = strings.ReplaceAll(between, "mermaid", "")
	// if between != "" {
	// 	m := mermaid.New()
	// 	pngUrl := m.RenderAsPng(between)
	// 	afterContent := utils.AfterString("```", chatMsg)
	// 	retMsg = fmt.Sprintf("![](%s)\n%s", pngUrl, afterContent)
	// }
	// xxxxx ```mermaid\n xxxxxx \n```
	return
}
