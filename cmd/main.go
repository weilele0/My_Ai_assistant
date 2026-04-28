package main

import (
	"My_AI_Assistant/internal/config"
	"My_AI_Assistant/internal/db"
	"My_AI_Assistant/internal/handler"
	"My_AI_Assistant/internal/model"
	"My_AI_Assistant/internal/repository"
	"My_AI_Assistant/internal/router"
	"My_AI_Assistant/internal/service"
	"fmt"
)

func main() {
	fmt.Println("🚀 My_AI_Assistant 正在启动...")

	// 初始化所有数据库连接
	db.InitChroma()               // Chroma 健康检查
	mysqlDB := db.InitMySQL()     // MySQL
	redisClient := db.InitRedis() // Redis

	// 自动迁移表结构（新增 Subscription / UserSubscription / UserQuota）
	mysqlDB.AutoMigrate(
		&model.User{},
		&model.Document{},
		&model.AITask{},
		&model.RAGHistory{},
		&model.Subscription{},
		&model.UserSubscription{},
		&model.UserQuota{},
	)

	// ── 初始化 Repository ──
	userRepo := repository.NewUserRepository(mysqlDB)
	documentRepo := repository.NewDocumentRepository(mysqlDB)
	taskRepo := repository.NewAITaskRepository(mysqlDB)
	ragHistoryRepo := repository.NewRAGHistoryRepository(mysqlDB)
	quotaRepo := repository.NewQuotaRepository(mysqlDB)
	subRepo := repository.NewSubscriptionRepository(mysqlDB)

	// 初始化默认套餐（幂等，已存在则跳过）
	subRepo.SeedDefaultPlans()

	// ── 初始化 Service ──
	userService := service.NewUserService(userRepo)
	aiService := service.NewAIService(config.LoadConfig().DeepSeekAPIKey)
	embeddingService := service.NewEmbeddingService(config.LoadConfig().QwenAPIKey)
	chromaService := service.NewChromaService(config.LoadConfig().ChromaURL)
	quotaService := service.NewQuotaService(quotaRepo, subRepo)

	documentService := service.NewDocumentService(documentRepo, embeddingService, chromaService)
	taskService := service.NewTaskService(aiService, taskRepo, redisClient)

	// 获取与 DocumentService 相同的 collectionID，确保读写同一个集合
	ragCollectionID, err := chromaService.GetOrCreateCollection("documents_v2")
	if err != nil {
		panic("RAG Chroma 集合初始化失败: " + err.Error())
	}
	ragService := service.NewRAGService(embeddingService, chromaService, aiService, ragHistoryRepo, ragCollectionID)

	// ── 初始化 Handler ──
	userHandler := handler.NewUserHandler(userService)
	aiHandler := handler.NewAIHandler(aiService)
	documentHandler := handler.NewDocumentHandler(documentService)
	taskHandler := handler.NewTaskHandler(taskService)
	ragHandler := handler.NewRAGHandler(ragService)
	quotaHandler := handler.NewQuotaHandler(quotaService)
	adminHandler := handler.NewAdminHandler(userRepo, quotaService, quotaRepo, taskRepo, ragHistoryRepo)

	// 设置路由
	r := router.SetupRouter(userHandler, aiHandler, documentHandler, taskHandler, ragHandler, quotaHandler, adminHandler)

	// 启动 Redis Stream Worker
	go taskService.StartTaskWorker()

	fmt.Println("\n✅ My_AI_Assistant 服务启动成功！")
	fmt.Println("📍 访问地址: http://localhost:8080")
	r.Run(":8080")
}
