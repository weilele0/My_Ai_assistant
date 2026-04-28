package model

import "time"

// RAGHistory RAG 生成历史记录
type RAGHistory struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	UserID     uint      `gorm:"index;not null" json:"user_id"`
	Question   string    `gorm:"type:text;not null" json:"question"`   // 用户问题
	Answer     string    `gorm:"type:text" json:"answer"`              // AI 回答
	References string    `gorm:"type:text" json:"-"`                   // 引用的文档片段（JSON 存储，不直接返回）
	ReferencesArr []string `gorm:"-" json:"references"`                // 转换后的引用列表
	CreatedAt  time.Time `json:"created_at"`
}
