package repository

import (
	"My_AI_Assistant/internal/model"

	"gorm.io/gorm"
)

// DocumentRepository 文档数据访问层
type DocumentRepository struct {
	db *gorm.DB
} //把数据库连接 *gorm.DB 包起来

// NewDocumentRepository 创建 DocumentRepository
func NewDocumentRepository(db *gorm.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
} //创建 DocumentRepository 实例 依赖注入

// Create 创建文档
func (r *DocumentRepository) Create(doc *model.Document) error {
	return r.db.Create(doc).Error //向 documents 表插入一条数据
}

// FindByUserID 根据用户ID查找文档列表
func (r *DocumentRepository) FindByUserID(userID uint) ([]model.Document, error) {
	var docs []model.Document //创建一个空结构体切片     //按照创建时间倒序                结果放入空结构体中
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&docs).Error
	return docs, err
}

// FindByID 根据ID查找文档（包含权限校验）
func (r *DocumentRepository) FindByID(id uint, userID uint) (*model.Document, error) {
	var doc model.Document
	//                      文章id       用户id                        第一条
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&doc).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// Update 更新文档
func (r *DocumentRepository) Update(doc *model.Document) error {
	return r.db.Save(doc).Error
}

// Delete 删除文档
func (r *DocumentRepository) Delete(id uint, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Document{}).Error
}
