package handler

import (
	"net/http"

	"My_AI_Assistant/internal/model"
	"My_AI_Assistant/internal/service"

	"github.com/gin-gonic/gin"
)

type AIHandler struct {
	aiService *service.AIService
}

// 创建并初始化AIHandler实例  传入已经初始化的ai服务实例  返回创建好的处理器实例
func NewAIHandler(aiService *service.AIService) *AIHandler {

	return &AIHandler{aiService: aiService}
}

// Generate 处理主题生成请求
func (h *AIHandler) Generate(c *gin.Context) {
	var req model.GenerateRequest //请求结构体实例化//用来接收前端请求
	/*Gin 上下文指针
	作用：获取请求数据（JSON/Query/Param）
	绑定参数、校验
	返回 JSON 响应
	存储请求生命周期数据*/
	//绑定并校验前端JSON参数
	// c.ShouldBindJSON(&req) Gin 框架提供的 JSON 参数绑定 + 校验二合一方法
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	///h.aiService.GenerateContent(...)处理器实例调用业务层，执行 AI 内容生成
	content, err := h.aiService.GenerateContent(req.Topic) //用户需要生成文章的主题
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI 生成失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "生成成功",
		"data": gin.H{
			"content": content,
		},
	})
}
