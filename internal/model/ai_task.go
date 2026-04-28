package model

import "time"

type AITask struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	UserID      uint       `gorm:"index;not null" json:"user_id"`
	TaskType    string     `gorm:"size:50;not null" json:"task_type"` // generate / improve / rag
	InputText   string     `gorm:"type:text" json:"input_text"`
	OutputText  string     `gorm:"type:text" json:"output_text"`
	Status      string     `gorm:"size:20;default:'pending'" json:"status"` // pending / processing / completed / failed
	ErrorMsg    string     `gorm:"type:text" json:"error_msg"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
}
