package handlers

import "github.com/xbmlz/chatgpt-plus-dingtalk/pkg/dingbot"

func HandlerHelp(ding *dingbot.DingBot) {
	content := `
	### 🤖 需要帮助吗？

	**我是卫博士，一款基于ChatGPT技术的智能聊天机器人！**

	🖼️ 生成图片👉 文本回复 *图片+空格+描述*

	🐳 流程图 👉 文本回复 *流程图+空格+描述*

	☘️ 帮助 👉 文本回复 *帮助*
	`
	ding.SendMessage(dingbot.MSG_MD, content)
}
