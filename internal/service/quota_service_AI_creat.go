package service

import (
	"My_AI_Assistant/internal/model"
	"My_AI_Assistant/internal/repository"
	"errors"
	"fmt"
	"time"
)

// QuotaService 用户额度管理服务
type QuotaService struct {
	quotaRepo *repository.QuotaRepository
	subRepo   *repository.SubscriptionRepository
}

// NewQuotaService 创建 QuotaService
func NewQuotaService(quotaRepo *repository.QuotaRepository, subRepo *repository.SubscriptionRepository) *QuotaService {
	return &QuotaService{
		quotaRepo: quotaRepo,
		subRepo:   subRepo,
	}
}

// QuotaStatus 额度状态（返回给前端或中间件）
type QuotaStatus struct {
	Plan         model.PlanType `json:"plan"`           // 当前套餐
	DailyAIQuota int            `json:"daily_ai_quota"` // 每日上限（-1=无限）
	DailyRAGQuota int           `json:"daily_rag_quota"`
	UsedAI       int            `json:"used_ai"`        // 今日已用
	UsedRAG      int            `json:"used_rag"`
	RemainAI     int            `json:"remain_ai"`      // 今日剩余（-1=无限）
	RemainRAG    int            `json:"remain_rag"`
	ExpiredAt    *time.Time     `json:"expired_at"` // 订阅到期时间（nil=永久）
}

// GetQuotaStatus 获取用户当前额度状态
func (s *QuotaService) GetQuotaStatus(userID uint) (*QuotaStatus, error) {
	// 1. 获取用户订阅
	us, err := s.subRepo.GetUserSubscription(userID)
	if err != nil {
		return nil, fmt.Errorf("获取订阅失败: %v", err)
	}

	// 2. 检查订阅是否过期
	plan := us.Plan
	if us.ExpiredAt != nil && us.ExpiredAt.Before(time.Now()) {
		plan = model.PlanFree // 过期降为免费版
	}

	// 3. 获取套餐限额
	sub, err := s.subRepo.GetPlanByName(plan)
	if err != nil {
		// 套餐表可能未初始化，使用默认免费限制
		sub = &model.Subscription{
			Plan:          model.PlanFree,
			DailyAIQuota:  5,
			DailyRAGQuota: 3,
		}
	}

	// 4. 获取今日用量
	quota, err := s.quotaRepo.GetOrCreate(userID)
	if err != nil {
		return nil, fmt.Errorf("获取用量失败: %v", err)
	}

	// 5. 计算剩余
	remainAI := -1
	remainRAG := -1
	if sub.DailyAIQuota >= 0 {
		remainAI = sub.DailyAIQuota - quota.UsedAI
		if remainAI < 0 {
			remainAI = 0
		}
	}
	if sub.DailyRAGQuota >= 0 {
		remainRAG = sub.DailyRAGQuota - quota.UsedRAG
		if remainRAG < 0 {
			remainRAG = 0
		}
	}

	return &QuotaStatus{
		Plan:          plan,
		DailyAIQuota:  sub.DailyAIQuota,
		DailyRAGQuota: sub.DailyRAGQuota,
		UsedAI:        quota.UsedAI,
		UsedRAG:       quota.UsedRAG,
		RemainAI:      remainAI,
		RemainRAG:     remainRAG,
		ExpiredAt:     us.ExpiredAt,
	}, nil
}

// CheckAndConsumeAI 校验 AI 额度是否足够，足够则扣减
// 返回 error 表示额度不足或系统错误
func (s *QuotaService) CheckAndConsumeAI(userID uint) error {
	status, err := s.GetQuotaStatus(userID)
	if err != nil {
		return err
	}
	// -1 表示无限制
	if status.DailyAIQuota >= 0 && status.RemainAI <= 0 {
		return errors.New(fmt.Sprintf("今日 AI 生成次数已达上限（%d 次），请升级套餐", status.DailyAIQuota))
	}
	return s.quotaRepo.IncrementAI(userID)
}

// CheckAndConsumeRAG 校验 RAG 额度是否足够，足够则扣减
func (s *QuotaService) CheckAndConsumeRAG(userID uint) error {
	status, err := s.GetQuotaStatus(userID)
	if err != nil {
		return err
	}
	if status.DailyRAGQuota >= 0 && status.RemainRAG <= 0 {
		return errors.New(fmt.Sprintf("今日 RAG 问答次数已达上限（%d 次），请升级套餐", status.DailyRAGQuota))
	}
	return s.quotaRepo.IncrementRAG(userID)
}

// SetUserPlan 管理员设置用户套餐
func (s *QuotaService) SetUserPlan(userID uint, plan model.PlanType, expiredAt *time.Time) error {
	// 校验套餐名合法
	validPlans := map[model.PlanType]bool{
		model.PlanFree:    true,
		model.PlanBasic:   true,
		model.PlanPro:     true,
		model.PlanUnlimit: true,
	}
	if !validPlans[plan] {
		return fmt.Errorf("无效的套餐类型: %s", plan)
	}
	return s.subRepo.UpsertUserSubscription(userID, plan, expiredAt)
}

// GetAllSubscriptions 管理员：获取所有用户订阅列表
func (s *QuotaService) GetAllSubscriptions() ([]model.UserSubscription, error) {
	return s.subRepo.GetAllUserSubscriptions()
}
