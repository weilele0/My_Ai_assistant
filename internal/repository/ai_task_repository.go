package repository

import (
	"My_AI_Assistant/internal/model"

	"gorm.io/gorm"
)

// AITaskRepository AI 任务数据访问层
type AITaskRepository struct {
	db *gorm.DB
}

// NewAITaskRepository 创建 AITaskRepository
func NewAITaskRepository(db *gorm.DB) *AITaskRepository {
	return &AITaskRepository{db: db}
}

// Create 创建任务
func (r *AITaskRepository) Create(task *model.AITask) error {
	return r.db.Create(task).Error
}

// FindByID 根据 ID 查找任务（带用户权限校验）
func (r *AITaskRepository) FindByID(id uint, userID uint) (*model.AITask, error) {
	var task model.AITask
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// Update 更新任务
func (r *AITaskRepository) Update(task *model.AITask) error {
	return r.db.Save(task).Error
}

// FindByUserID 获取用户的所有任务
func (r *AITaskRepository) FindByUserID(userID uint) ([]model.AITask, error) {
	var tasks []model.AITask
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&tasks).Error
	return tasks, err
}

// FindByIDForWorker 后台 worker 专用（不校验 userID）
func (r *AITaskRepository) FindByIDForWorker(id uint) (*model.AITask, error) {
	var task model.AITask
	err := r.db.First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}
