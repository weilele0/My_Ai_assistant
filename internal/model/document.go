package model

import "time"

// Document 文档模型
type Document struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"` // 所属用户ID
	Title     string    `gorm:"size:255;not null" json:"title"`
	FilePath  string    `gorm:"size:255" json:"file_path"` // 文件存储路径
	Content   string    `gorm:"type:text" json:"content"`  // 提取后的文本内容
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
