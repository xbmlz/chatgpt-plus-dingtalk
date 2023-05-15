package dingbot

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/fetch"
)

type MsgType string

const MSG_TEXT MsgType = "text"
const MSG_MD MsgType = "markdown"

type DingBot struct {
	Msg DingBotReceiveMessage
}

type DingBotText struct {
	Content string `json:"content"` // 消息文本
}

type DingBotMarkdown struct {
	Text  string `json:"text"`  // markdown 消息内容
	Title string `json:"title"` // markdown 消息标题
}

type DingBotLink struct {
	MessageUrl string `json:"messageUrl"` // 消息点击链接地址
	Title      string `json:"title"`      // 消息标题
	PicUrl     string `json:"picUrl"`     // 图片地址
	Text       string `json:"text"`       // 消息内容。如果太长只会部分展示
}

// See https://open.dingtalk.com/document/orgapp/the-application-robot-in-the-enterprise-sends-a-single-chat
type DingBotReceiveMessage struct {
	ConversationID string `json:"conversationId"` // 群ID
	AtUsers        []struct {
		DingtalkID string `json:"dingtalkId"` // 加密的发送者ID
	} `json:"atUsers"` // 被@人的信息
	ChatbotUserID             string      `json:"chatbotUserId"`             // 加密的发送者ID
	MsgID                     string      `json:"msgId"`                     // 消息ID
	SenderNick                string      `json:"senderNick"`                // 发送者昵称
	IsAdmin                   bool        `json:"isAdmin"`                   // 是否是管理员
	SenderStaffID             string      `json:"senderStaffId"`             // 企业内部群中@该机器人的成员userid
	SessionWebhookExpiredTime int64       `json:"sessionWebhookExpiredTime"` // 会话过期时间戳，单位毫秒
	CreateAt                  int64       `json:"createAt"`                  // 消息创建时间戳，单位毫秒
	ConversationType          string      `json:"conversationType"`          // 会话类型，1表示私聊 2表示群聊
	SenderID                  string      `json:"senderId"`                  // @该机器人的成员的加密ID
	ConversationTitle         string      `json:"conversationTitle"`         // 会话标题
	IsInAtList                bool        `json:"isInAtList"`                // 是否在@列表中
	SessionWebhook            string      `json:"sessionWebhook"`            // 当前会话的Webhook地址
	Text                      DingBotText `json:"text"`                      // 消息文本
	RobotCode                 string      `json:"robotCode"`                 // 机器人code 自定义机器人默认为normal
	MsgType                   string      `json:"msgtype"`                   // 消息类型 目前只支持text
}

type DingBotSendMessage struct {
	MsgType  string          `json:"msgtype"`            // 消息类型 支持 text, markdown
	Markdown DingBotMarkdown `json:"markdown,omitempty"` // markdown 消息内容
	Text     DingBotText     `json:"text,omitempty"`     // 消息文本
	At       DingBotSendAt   `json:"at"`                 // @人的信息
}

type DingBotSendAt struct {
	AtUserIds []string `json:"atUserIds"`
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

type DingBotResponse struct {
	ErrCode int    `json:"errcode"` // 错误码    0
	ErrMsg  string `json:"errmsg"`  // 错误信息 "ok"
}

func New(msg DingBotReceiveMessage) *DingBot {
	return &DingBot{
		Msg: msg,
	}
}

func (d *DingBot) SendMessage(mtype MsgType, msg string) (err error) {
	if d.Msg.SenderStaffID == "" {
		msg = fmt.Sprintf("%s\n\n@%s", msg, d.Msg.SenderNick)
	}
	var send DingBotSendMessage
	switch mtype {
	case MSG_TEXT:
		send = DingBotSendMessage{
			MsgType: string(MSG_TEXT),
			Text: DingBotText{
				Content: msg,
			},
		}
	case MSG_MD:
		send = DingBotSendMessage{
			MsgType: string(MSG_MD),
			Markdown: DingBotMarkdown{
				Text:  msg,
				Title: msg,
			},
		}
	default:
		return errors.New("不支持的消息类型")
	}
	if d.Msg.ConversationType == "2" {
		send.At = DingBotSendAt{
			AtUserIds: []string{d.Msg.SenderID},
		}
	}
	data, err := json.Marshal(send)
	if err != nil {
		return err
	}
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "*/*",
	}
	resp, err := fetch.POST(d.Msg.SessionWebhook, headers, data)
	if err != nil {
		return err
	}
	response := DingBotResponse{}
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return err
	}
	if response.ErrCode != 0 {
		return errors.New(response.ErrMsg)
	}
	return nil
}
