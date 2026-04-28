package model

import "time"

// UserQuota 用户每日额度使用情况
// 每个用户每天一条记录，记录当天已用次数
type UserQuota struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	UserID       uint      `gorm:"index;not null" json:"user_id"`
	Date         string    `gorm:"size:10;index;not null" json:"date"` // 日期，格式 2006-01-02
	UsedAI       int       `gorm:"not null;default:0" json:"used_ai"`  // 当日已用 AI 生成次数
	UsedRAG      int       `gorm:"not null;default:0" json:"used_rag"` // 当日已用 RAG 问答次数
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
