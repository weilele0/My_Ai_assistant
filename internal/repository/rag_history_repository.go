package repository

import (
	"My_AI_Assistant/internal/model"

	"gorm.io/gorm"
)

// RAGHistoryRepository RAG 历史记录数据访问层
type RAGHistoryRepository struct {
	db *gorm.DB
}

// NewRAGHistoryRepository 创建 RAGHistoryRepository
func NewRAGHistoryRepository(db *gorm.DB) *RAGHistoryRepository {
	return &RAGHistoryRepository{db: db}
}

// Create 创建 RAG 历史记录，并自动裁剪至保留最近 limit 条
func (r *RAGHistoryRepository) Create(history *model.RAGHistory, keepLimit int) error {
	err := r.db.Create(history).Error
	if err != nil {
		return err
	}

	// 裁剪：只保留该用户最近 keepLimit 条
	r.db.Where("user_id = ? AND id NOT IN (SELECT id FROM rag_histories WHERE user_id = ? ORDER BY created_at DESC LIMIT ?)",
		history.UserID, history.UserID, keepLimit).
		Delete(&model.RAGHistory{})

	return nil
}

// FindByUserID 获取用户所有 RAG 历史（倒序，最新的在前）
func (r *RAGHistoryRepository) FindByUserID(userID uint) ([]model.RAGHistory, error) {
	var histories []model.RAGHistory
	err := r.db.Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&histories).Error
	return histories, err
}

// FindByID 根据 ID 查找（带用户权限校验）
func (r *RAGHistoryRepository) FindByID(id uint, userID uint) (*model.RAGHistory, error) {
	var history model.RAGHistory
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&history).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

// Delete 删除指定记录（带用户权限校验）
func (r *RAGHistoryRepository) Delete(id uint, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.RAGHistory{}).Error
}

// ClearAll 清空用户所有 RAG 历史
func (r *RAGHistoryRepository) ClearAll(userID uint) error {
	//调用数据库  用户id           删除这个                     指定表名
	return r.db.Where("user_id = ?", userID).Delete(&model.RAGHistory{}).Error
}
