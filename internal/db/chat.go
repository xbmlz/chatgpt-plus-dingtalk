package db

import (
	"gorm.io/gorm"
)

type ChatMode uint

const QUESTION_MODE ChatMode = 1
const ANSWER_MODE ChatMode = 2

type Chat struct {
	gorm.Model
	SenderID               string   `gorm:"type:varchar(128);not null;comment:'发送者ID'" json:"sender_id"`
	SenderNick             string   `gorm:"type:varchar(50);not null;comment:'发送者昵称'" json:"sender_nick"`
	DingTalkConversationID string   `gorm:"type:varchar(128);not null;comment:'钉钉会话ID'" json:"ding_talk_conversation_id"`
	MessageID              string   `gorm:"type:varchar(128);not null;comment:'消息ID'" json:"message_id"`
	ConversationID         string   `gorm:"type:varchar(128);not null;comment:'会话ID'" json:"conversation_id"`
	ConversationMode       ChatMode `gorm:"type:int;not null;comment:'消息模式1问题2回答'" json:"conversation_mode"`
	ConversationType       string   `gorm:"type:int;not null;comment:'消息类型1私聊2群聊'" json:"conversation_type"`
	Content                string   `gorm:"type:varchar(255);not null;comment:'消息内容'" json:"content"`
}

type ChatListParams struct {
	SenderID string `form:"sender_id" json:"sender_id"`
}

// Create
func Save(c *Chat) error {
	return DB.Create(&c).Error
}

// Find
func FindOne(filter map[string]interface{}, chat *Chat) error {
	return DB.Where(filter).Last(&chat).Error
}

// Delete
func DeleteByDingTalkConversationID(dingTalkConversation_id string) error {
	return DB.Where("ding_talk_conversation_id = ?", dingTalkConversation_id).Delete(&Chat{}).Error
}

// FindAllDingTalkConversationId find all ding talk conversation id
func FindAllDingTalkConversationId() ([]string, error) {
	var chat []Chat
	var conversationID []string
	err := DB.Distinct("ding_talk_conversation_id").Find(&chat).Error
	if err != nil {
		return conversationID, err
	}
	for _, v := range chat {
		conversationID = append(conversationID, v.DingTalkConversationID)
	}
	return conversationID, nil
}

// FindGptConversationId find gpt conversation id
func FindGptConversationId(dingTalkConversationId string) (string, error) {
	var chat Chat
	err := DB.Where("ding_talk_conversation_id = ? and conversation_mode = 2", dingTalkConversationId).Last(&chat).Error
	if err != nil {
		return "", err
	}
	return chat.ConversationID, nil
}
