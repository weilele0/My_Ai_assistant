package middleware

import (
	"My_AI_Assistant/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// QuotaCheckAI 检查用户 AI 生成额度中间件
// 路由使用方式：router.Use(middleware.QuotaCheckAI(quotaService))
func QuotaCheckAI(quotaService *service.QuotaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
			c.Abort()
			return
		}

		if err := quotaService.CheckAndConsumeAI(userID.(uint)); err != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":  429,
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// QuotaCheckRAG 检查用户 RAG 问答额度中间件
func QuotaCheckRAG(quotaService *service.QuotaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
			c.Abort()
			return
		}

		if err := quotaService.CheckAndConsumeRAG(userID.(uint)); err != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":  429,
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
