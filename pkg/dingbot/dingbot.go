package dingbot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type MsgType string

const TEXT MsgType = "text"
const MARKDOWN MsgType = "markdown"

type DingBot struct {
	SessionWebhook string // 当前会话的Webhook地址
}

type DingBotText struct {
	Content string `json:"content"` // 消息文本
}

type DingBotMarkdown struct {
	Text  string `json:"text"`  // markdown 消息内容
	Title string `json:"title"` // markdown 消息标题
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

func NewDingBot(dingbot DingBot) *DingBot {
	return &DingBot{
		SessionWebhook: dingbot.SessionWebhook,
	}
}

func (d *DingBot) SendMarkdownMessage(ctype, mtype, msg, atID, atNick string) (err error) {
	if atID == "" {
		msg = fmt.Sprintf("%s\n\n@%s", msg, atNick)
	}
	var send DingBotSendMessage
	switch mtype {
	case "text":
		send = DingBotSendMessage{
			MsgType: "markdown",
			Markdown: DingBotMarkdown{
				Text:  msg,
				Title: msg,
			},
		}
	case "markdown":
		send = DingBotSendMessage{
			MsgType: "markdown",
			Text: DingBotText{
				Content: msg,
			},
		}
	}
	if ctype == "2" {
		send.At = DingBotSendAt{
			AtUserIds: []string{atID},
		}
	}
	data, err := json.Marshal(send)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", d.SessionWebhook, bytes.NewBuffer(data))

	if err != nil {
		return err
	}

	req.Header.Add("Accept", "*/*")
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	response := DingBotResponse{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}
	if response.ErrCode != 0 {
		return errors.New(response.ErrMsg)
	}
	return nil
}
