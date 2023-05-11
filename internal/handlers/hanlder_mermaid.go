package handlers

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/xbmlz/chatgpt-plus-dingtalk/internal/config"
	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/dingbot"
	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/mermaid"
	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/utils"
)

func HandlerMermaid(msg dingbot.DingBotReceiveMessage) (retMsg string) {
	chatMsg := HandlerMessage(msg)
	between := utils.ExtractStringBetween("```mermaid", "```", chatMsg)
	if between != "" {
		ctx := context.Background()
		ret, _ := mermaid.NewRenderEngine(ctx, `mermaid.initialize({'theme': 'base', 'themeVariables': { 'primaryColor': '#1473e6'}});`)
		defer ret.Cancel()
		resultBytes, box, err := ret.RenderAsPng(between)
		if err != nil {
			retMsg = err.Error()
			return
		}
		if box.Width < 1 || box.Height < 1 {
			retMsg = fmt.Sprintf("Render() got empty image = w:%d, h:%d)", box.Width, box.Height)
			return
		}
		fileName := time.Now().Format("20060102150405") + ".png"
		filePath := filepath.Join("./data/images/", fileName)
		err = utils.ImageSave(resultBytes, filePath)
		if err != nil {
			retMsg = err.Error()
			return
		}
		fileUrl := fmt.Sprintf("%s/images/%s", config.Instance.ServerUrl, fileName)
		afterContent := utils.AfterString("```", chatMsg)
		// pngStr := fmt.Sprintf("![](%s)\n%s", fileUrl, afterContent)
		// chatMsg = strings.ReplaceAll(chatMsg, between, pngStr)
		retMsg = fmt.Sprintf("![](%s)\n%s", fileUrl, afterContent)
	}
	// xxxxx ```mermaid\n xxxxxx \n```
	return
}
