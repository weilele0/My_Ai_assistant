package model

import "time"

type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Username  string    `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email     string    `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"-"`               // 不返回给前端
	IsAdmin   bool      `gorm:"not null;default:false" json:"is_admin"`   // 是否为管理员
	CreatedAt time.Time `json:"created_at"`
}
