package handler

import (
	"My_AI_Assistant/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// QuotaHandler 用户额度控制器
type QuotaHandler struct {
	quotaService *service.QuotaService
}

// NewQuotaHandler 创建 QuotaHandler
func NewQuotaHandler(quotaService *service.QuotaService) *QuotaHandler {
	return &QuotaHandler{quotaService: quotaService}
}

// QuotaService 暴露内部 quotaService，供 router 注册中间件时使用
func (h *QuotaHandler) QuotaService() *service.QuotaService {
	return h.quotaService
}

// GetMyQuota 获取当前用户的额度状态
// GET /api/v1/quota/me
func (h *QuotaHandler) GetMyQuota(c *gin.Context) {
	userID, _ := c.Get("user_id")

	status, err := h.quotaService.GetQuotaStatus(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取额度失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": status,
	})
}
