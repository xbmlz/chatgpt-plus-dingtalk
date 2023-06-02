package chatgpt

import (
	"bufio"
	"bytes"
	"encoding/json"
	"strings"

	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/fetch"
	"github.com/xbmlz/chatgpt-plus-dingtalk/pkg/logger"
)

type ChatGPT struct {
	BaseUrl     string
	AccessToken string
}

type CompletionMessageContent struct {
	ContentType string   `json:"content_type"` // 消息类型 目前只支持text
	Parts       []string `json:"parts"`        // 消息文本
}

// request
type CompletionRequestMessage struct {
	ID      string                   `json:"id"`      // 消息ID
	Role    string                   `json:"role"`    // 消息角色 user system assistant
	Content CompletionMessageContent `json:"content"` // 消息内容
}

type CompletionRequest struct {
	Action          string                     `json:"action"`                      // 消息类型 目前只支持next
	Messages        []CompletionRequestMessage `json:"messages"`                    // 消息内容
	Model           string                     `json:"model"`                       // 消息模型 text-davinci-002-render-sha
	ParentMessageID string                     `json:"parent_message_id,omitempty"` // 父消息ID
	ConversationID  string                     `json:"conversation_id,omitempty"`   // 会话ID
}

// response
type CompletionResponseMessageAuthor struct {
	Role     string                 `json:"role"`      // 消息角色 user system assistant
	Name     string                 `json:"name"`      // 角色名称
	MetaData map[string]interface{} `json:"meta_data"` // meta_data
}

type CompletionResponseMessage struct {
	ID         string                          `json:"id"`          // 消息ID
	Author     CompletionResponseMessageAuthor `json:"author"`      // 消息角色 user system assistant
	CreateTime float64                         `json:"create_time"` // 创建时间
	UpdateTime float64                         `json:"update_time"` // 更新时间
	Content    CompletionMessageContent        `json:"content"`     // 消息内容
	EndTurn    bool                            `json:"end_turn"`    // 是否结束会话
	Weight     float64                         `json:"weight"`      // 权重
	MetaData   map[string]interface{}          `json:"meta_data"`   // meta_data
	Recipient  string                          `json:"recipient"`   // 接收者
}

type CompletionResponse struct {
	ConversationID string                    `json:"conversation_id"` // 会话ID
	Error          string                    `json:"error"`           // 错误信息
	Message        CompletionResponseMessage `json:"message"`         // 消息内容
}

func New(options ChatGPT) *ChatGPT {
	return &ChatGPT{
		AccessToken: options.AccessToken,
		BaseUrl:     options.BaseUrl,
	}
}

func (c *ChatGPT) CreateCompletion(param CompletionRequest) (res CompletionResponse, err error) {
	data, err := json.Marshal(param)
	if err != nil {
		return res, err
	}
	headers := map[string]string{
		"Accept":        "text/event-stream",
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + c.AccessToken,
	}
	resp, err := fetch.POST(c.BaseUrl+"/conversation", headers, data)
	if err != nil {
		return res, err
	}
	scanner := bufio.NewScanner(bytes.NewReader(resp))
	for scanner.Scan() {
		line := scanner.Text()
		// 处理每一行数据
		if strings.HasPrefix(line, "data:") {
			line = strings.TrimPrefix(line, "data: ")
			if line == "[DONE]" {
				break
			}
			err := json.Unmarshal([]byte(line), &res)
			if err != nil {
				logger.Error(err)
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Error(err)
		return res, err
	}

	return res, nil
}

func (c *ChatGPT) DeleteConversation(id string) (err error) {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + c.AccessToken,
	}
	_, err = fetch.DELETE(c.BaseUrl+"/conversation/"+id, headers)
	if err != nil {
		return
	}
	return nil
}
