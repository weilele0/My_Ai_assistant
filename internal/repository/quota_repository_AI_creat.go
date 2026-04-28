package repository

import (
	"My_AI_Assistant/internal/model"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// QuotaRepository 用户每日额度数据访问层
type QuotaRepository struct {
	db *gorm.DB
}

// NewQuotaRepository 创建 QuotaRepository
func NewQuotaRepository(db *gorm.DB) *QuotaRepository {
	return &QuotaRepository{db: db}
}

// today 返回今天的日期字符串（2006-01-02）
func today() string {
	return time.Now().Format("2006-01-02")
}

// GetOrCreate 获取用户今日额度记录，不存在则创建
func (r *QuotaRepository) GetOrCreate(userID uint) (*model.UserQuota, error) {
	var quota model.UserQuota
	date := today()
	// Upsert：如果不存在则插入，如果存在则什么都不做
	result := r.db.Where(model.UserQuota{UserID: userID, Date: date}).
		Attrs(model.UserQuota{UsedAI: 0, UsedRAG: 0}).
		FirstOrCreate(&quota)
	return &quota, result.Error
}

// IncrementAI 将用户今日 AI 使用次数 +1
func (r *QuotaRepository) IncrementAI(userID uint) error {
	date := today()
	// 先确保记录存在
	r.db.Where(model.UserQuota{UserID: userID, Date: date}).
		Attrs(model.UserQuota{UsedAI: 0, UsedRAG: 0}).
		FirstOrCreate(&model.UserQuota{})

	return r.db.Model(&model.UserQuota{}).
		Where("user_id = ? AND date = ?", userID, date).
		UpdateColumn("used_ai", gorm.Expr("used_ai + 1")).Error
}

// IncrementRAG 将用户今日 RAG 使用次数 +1
func (r *QuotaRepository) IncrementRAG(userID uint) error {
	date := today()
	r.db.Where(model.UserQuota{UserID: userID, Date: date}).
		Attrs(model.UserQuota{UsedAI: 0, UsedRAG: 0}).
		FirstOrCreate(&model.UserQuota{})

	return r.db.Model(&model.UserQuota{}).
		Where("user_id = ? AND date = ?", userID, date).
		UpdateColumn("used_rag", gorm.Expr("used_rag + 1")).Error
}

// FindAllByUserID 查询用户所有额度历史
func (r *QuotaRepository) FindAllByUserID(userID uint) ([]model.UserQuota, error) {
	var quotas []model.UserQuota
	err := r.db.Where("user_id = ?", userID).Order("date desc").Find(&quotas).Error
	return quotas, err
}

// AdminGetAllQuotaToday 管理员：查询今日所有用户的额度使用情况
func (r *QuotaRepository) AdminGetAllQuotaToday() ([]model.UserQuota, error) {
	var quotas []model.UserQuota
	err := r.db.Where("date = ?", today()).Order("used_ai + used_rag desc").Find(&quotas).Error
	return quotas, err
}

// 确保 clause 包被使用（避免 unused import）
var _ = clause.OnConflict{}
