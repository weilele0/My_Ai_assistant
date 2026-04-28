package model

import "time"

// PlanType 套餐类型
type PlanType string

const (
	PlanFree    PlanType = "free"    // 免费版
	PlanBasic   PlanType = "basic"   // 基础版
	PlanPro     PlanType = "pro"     // 专业版
	PlanUnlimit PlanType = "unlimit" // 无限制（管理员手动赋予）
)

// Subscription 订阅套餐定义表（套餐模板）
// 记录每种套餐每日可调用的 AI 次数上限
type Subscription struct {
	ID            uint     `gorm:"primarykey" json:"id"`
	Plan          PlanType `gorm:"size:20;uniqueIndex;not null" json:"plan"`           // 套餐名称
	DailyAIQuota  int      `gorm:"not null;default:10" json:"daily_ai_quota"`          // 每日 AI 生成次数上限（-1=无限）
	DailyRAGQuota int      `gorm:"not null;default:5" json:"daily_rag_quota"`          // 每日 RAG 问答次数上限（-1=无限）
	Description   string   `gorm:"size:255" json:"description"`                        // 套餐说明
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// UserSubscription 用户订阅记录（用户当前生效的套餐）
type UserSubscription struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex;not null" json:"user_id"` // 每人只有一条生效记录
	Plan      PlanType  `gorm:"size:20;not null;default:'free'" json:"plan"`
	ExpiredAt *time.Time `json:"expired_at"` // nil = 永久有效
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
