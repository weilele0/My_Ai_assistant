package repository

import (
	"My_AI_Assistant/internal/model"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SubscriptionRepository 订阅套餐数据访问层
type SubscriptionRepository struct {
	db *gorm.DB
}

// NewSubscriptionRepository 创建 SubscriptionRepository
func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

// GetPlanByName 根据套餐名查询套餐信息
func (r *SubscriptionRepository) GetPlanByName(plan model.PlanType) (*model.Subscription, error) {
	var sub model.Subscription
	err := r.db.Where("plan = ?", plan).First(&sub).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

// GetUserSubscription 获取用户当前订阅
func (r *SubscriptionRepository) GetUserSubscription(userID uint) (*model.UserSubscription, error) {
	var us model.UserSubscription
	err := r.db.Where("user_id = ?", userID).First(&us).Error
	if err != nil {
		// 用户没有订阅记录，返回默认免费套餐
		if err == gorm.ErrRecordNotFound {
			return &model.UserSubscription{
				UserID: userID,
				Plan:   model.PlanFree,
			}, nil
		}
		return nil, err
	}
	return &us, nil
}

// UpsertUserSubscription 创建或更新用户订阅（管理员操作）
func (r *SubscriptionRepository) UpsertUserSubscription(userID uint, plan model.PlanType, expiredAt *time.Time) error {
	us := model.UserSubscription{
		UserID:    userID,
		Plan:      plan,
		ExpiredAt: expiredAt,
	}
	// 如果存在则更新 plan 和 expired_at，不存在则插入
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"plan", "expired_at", "updated_at"}),
	}).Create(&us).Error
}

// GetAllUserSubscriptions 管理员：查询所有用户订阅
func (r *SubscriptionRepository) GetAllUserSubscriptions() ([]model.UserSubscription, error) {
	var subs []model.UserSubscription
	err := r.db.Order("updated_at desc").Find(&subs).Error
	return subs, err
}

// SeedDefaultPlans 初始化默认套餐（幂等）
func (r *SubscriptionRepository) SeedDefaultPlans() {
	plans := []model.Subscription{
		{Plan: model.PlanFree, DailyAIQuota: 5, DailyRAGQuota: 3, Description: "免费版：每日 AI 生成 5 次，RAG 问答 3 次"},
		{Plan: model.PlanBasic, DailyAIQuota: 30, DailyRAGQuota: 20, Description: "基础版：每日 AI 生成 30 次，RAG 问答 20 次"},
		{Plan: model.PlanPro, DailyAIQuota: 100, DailyRAGQuota: 60, Description: "专业版：每日 AI 生成 100 次，RAG 问答 60 次"},
		{Plan: model.PlanUnlimit, DailyAIQuota: -1, DailyRAGQuota: -1, Description: "无限制版：不限次数（仅管理员可分配）"},
	}
	for _, p := range plans {
		r.db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "plan"}},
			DoNothing: true, // 已存在则跳过，不覆盖管理员的修改
		}).Create(&p)
	}
}
