package handler

import (
	"My_AI_Assistant/internal/model"
	"My_AI_Assistant/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RAGHandler struct {
	ragService *service.RAGService
}

func NewRAGHandler(ragService *service.RAGService) *RAGHandler {
	return &RAGHandler{ragService: ragService}
}

// GenerateWithRAG 处理 RAG 生成请求
func (h *RAGHandler) GenerateWithRAG(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req model.RAGRequest
	//解析json请求到req结构体里
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//调用服务层，去包装请求，构建promat
	result, err := h.ragService.GenerateWithRAG(userID.(uint), req.Question, req.TopK)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": result,
	})
}

// GetHistory 获取 RAG 历史记录列表
func (h *RAGHandler) GetHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")

	histories, err := h.ragService.GetHistory(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取历史记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": histories,
	})
}

// GetHistoryDetail 获取单条 RAG 历史详情
func (h *RAGHandler) GetHistoryDetail(c *gin.Context) {
	userID, _ := c.Get("user_id")
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	history, err := h.ragService.GetHistoryDetail(uint(id), userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "历史记录不存在或无权限"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": history,
	})
}

// DeleteHistory 删除单条 RAG 历史
func (h *RAGHandler) DeleteHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	if err := h.ragService.DeleteHistory(uint(id), userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}

// ClearHistory 清空所有 RAG 历史
func (h *RAGHandler) ClearHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")

	if err := h.ragService.ClearHistory(userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "清空失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "已清空所有历史记录",
	})
}
