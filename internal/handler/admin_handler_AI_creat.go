package handler

import (
	"My_AI_Assistant/internal/model"
	"My_AI_Assistant/internal/repository"
	"My_AI_Assistant/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// AdminHandler 管理员后台控制器
type AdminHandler struct {
	userRepo     *repository.UserRepository
	quotaService *service.QuotaService
	quotaRepo    *repository.QuotaRepository
	taskRepo     *repository.AITaskRepository
	ragHistRepo  *repository.RAGHistoryRepository
}

// NewAdminHandler 创建 AdminHandler
func NewAdminHandler(
	userRepo *repository.UserRepository,
	quotaService *service.QuotaService,
	quotaRepo *repository.QuotaRepository,
	taskRepo *repository.AITaskRepository,
	ragHistRepo *repository.RAGHistoryRepository,
) *AdminHandler {
	return &AdminHandler{
		userRepo:     userRepo,
		quotaService: quotaService,
		quotaRepo:    quotaRepo,
		taskRepo:     taskRepo,
		ragHistRepo:  ragHistRepo,
	}
}

// GetUserList 管理员：获取用户列表
// GET /api/v1/admin/users
func (h *AdminHandler) GetUserList(c *gin.Context) {
	users, err := h.userRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": users,
	})
}

// GetUserDetail 管理员：获取单个用户详情 + 额度信息
// GET /api/v1/admin/users/:id
func (h *AdminHandler) GetUserDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	user, err := h.userRepo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	quota, err := h.quotaService.GetQuotaStatus(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取额度失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"user":  user,
			"quota": quota,
		},
	})
}

// GetStats 管理员：查看系统统计数据
// GET /api/v1/admin/stats
func (h *AdminHandler) GetStats(c *gin.Context) {
	// 用户总数
	userCount, _ := h.userRepo.Count()

	// 今日额度使用情况
	todayQuotas, _ := h.quotaRepo.AdminGetAllQuotaToday()

	// 统计今日 AI 总调用次数和 RAG 总调用次数
	var totalAI, totalRAG int
	for _, q := range todayQuotas {
		totalAI += q.UsedAI
		totalRAG += q.UsedRAG
	}

	// 所有用户订阅分布
	subs, _ := h.quotaService.GetAllSubscriptions()
	planDist := map[string]int{
		string(model.PlanFree):    0,
		string(model.PlanBasic):   0,
		string(model.PlanPro):     0,
		string(model.PlanUnlimit): 0,
	}
	for _, s := range subs {
		planDist[string(s.Plan)]++
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total_users":           userCount,
			"today_active_users":    len(todayQuotas),
			"today_ai_calls":        totalAI,
			"today_rag_calls":       totalRAG,
			"plan_distribution":     planDist,
			"today_quota_details":   todayQuotas,
		},
	})
}

// SetUserPlan 管理员：设置用户套餐
// POST /api/v1/admin/users/:id/plan
// Body: { "plan": "basic", "expired_at": "2026-12-31T23:59:59Z" }  expired_at 可选
func (h *AdminHandler) SetUserPlan(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	var req struct {
		Plan      model.PlanType `json:"plan" binding:"required"`
		ExpiredAt *time.Time     `json:"expired_at"` // 可选，nil = 永久
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.quotaService.SetUserPlan(uint(id), req.Plan, req.ExpiredAt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "套餐设置成功",
	})
}

// GetAllSubscriptions 管理员：查看所有用户订阅
// GET /api/v1/admin/subscriptions
func (h *AdminHandler) GetAllSubscriptions(c *gin.Context) {
	subs, err := h.quotaService.GetAllSubscriptions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取订阅列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": subs,
	})
}

// GetTodayQuotas 管理员：查看今日所有用户的额度使用情况
// GET /api/v1/admin/quotas/today
func (h *AdminHandler) GetTodayQuotas(c *gin.Context) {
	quotas, err := h.quotaRepo.AdminGetAllQuotaToday()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取额度数据失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": quotas,
	})
}
