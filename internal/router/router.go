package router

import (
	"My_AI_Assistant/internal/handler"
	"My_AI_Assistant/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	userHandler *handler.UserHandler,
	aiHandler *handler.AIHandler,
	documentHandler *handler.DocumentHandler,
	taskHandler *handler.TaskHandler,
	ragHandler *handler.RAGHandler,
	quotaHandler *handler.QuotaHandler,
	adminHandler *handler.AdminHandler,
) *gin.Engine {
	r := gin.Default()

	// 静态文件服务 — 前端页面
	// 必须放在 api.Group 之前，且不能使用 /*filepath 通配符，否则与 /api/v1/* 冲突
	r.GET("/", func(c *gin.Context) {
		c.File("./frontend/index.html")
	})
	r.Static("/css", "./frontend/css")
	r.Static("/js", "./frontend/js")

	// SPA 兜底：所有未匹配的 GET 请求返回 index.html（必须放在 api.Group 之后）
	r.NoRoute(NoRoute)

	api := r.Group("/api/v1")
	{
		// ── 用户接口（公开）──
		user := api.Group("/user")
		{
			user.POST("/register", userHandler.Register)
			user.POST("/login", userHandler.Login)
		}

		// ── AI 接口（需要登录 + 额度校验）──
		ai := api.Group("/ai")
		ai.Use(middleware.JWTAuth())
		{
			// 同步 AI 生成：消耗 AI 额度
			ai.POST("/generate", middleware.QuotaCheckAI(quotaHandler.QuotaService()), aiHandler.Generate)
			// RAG 问答：消耗 RAG 额度
			ai.POST("/rag-generate", middleware.QuotaCheckRAG(quotaHandler.QuotaService()), ragHandler.GenerateWithRAG)
			// RAG 历史记录（不消耗额度）
			ai.GET("/rag-history", ragHandler.GetHistory)
			ai.GET("/rag-history/:id", ragHandler.GetHistoryDetail)
			ai.DELETE("/rag-history/:id", ragHandler.DeleteHistory)
			ai.DELETE("/rag-history", ragHandler.ClearHistory)
		}

		// ── 额度接口（需要登录）──
		quota := api.Group("/quota")
		quota.Use(middleware.JWTAuth())
		{
			quota.GET("/me", quotaHandler.GetMyQuota) // 查看自己的额度
		}

		// ── 文档接口（需要登录）──
		doc := api.Group("/documents")
		doc.Use(middleware.JWTAuth())
		{
			doc.POST("", documentHandler.UploadDocument)
			doc.GET("", documentHandler.GetUserDocuments)
			doc.GET("/:id", documentHandler.GetDocumentDetail)
			doc.DELETE("/:id", documentHandler.DeleteDocument)
		}

		// ── 任务接口（需要登录）──
		task := api.Group("/tasks")
		task.Use(middleware.JWTAuth())
		{
			task.POST("/generate", taskHandler.CreateAsyncGenerateTask)
			task.GET("", taskHandler.GetUserTasks)
			task.GET("/:id", taskHandler.GetTaskStatus)
		}

		// ── 管理员后台（需要登录 + 管理员权限）──
		admin := api.Group("/admin")
		admin.Use(middleware.JWTAuth(), middleware.AdminAuth())
		{
			// 用户管理
			admin.GET("/users", adminHandler.GetUserList)              // 用户列表
			admin.GET("/users/:id", adminHandler.GetUserDetail)        // 用户详情 + 额度
			admin.POST("/users/:id/plan", adminHandler.SetUserPlan)    // 设置用户套餐

			// 统计数据
			admin.GET("/stats", adminHandler.GetStats)                 // 系统统计

			// 额度管理
			admin.GET("/quotas/today", adminHandler.GetTodayQuotas)    // 今日所有用户额度
			admin.GET("/subscriptions", adminHandler.GetAllSubscriptions) // 所有订阅信息
		}
	}

	return r
}

// NoRoute 兜底：所有未匹配的 GET 请求都返回 index.html（SPA 路由兜底）
func NoRoute(c *gin.Context) {
	c.File("./frontend/index.html")
}
