package handler

import (
	"My_AI_Assistant/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DocumentHandler 文档控制器
type DocumentHandler struct {
	documentService *service.DocumentService
} //定义结构体注入服务层

// NewDocumentHandler 创建 DocumentHandler
func NewDocumentHandler(documentService *service.DocumentService) *DocumentHandler {
	return &DocumentHandler{documentService: documentService}
} //实例化，传入service

// UploadDocument 处理文档上传请求
func (h *DocumentHandler) UploadDocument(c *gin.Context) {
	userID, _ := c.Get("user_id") // 从 JWT 中间件获取当前用户ID
	//绑定并校验前端的参数
	var req struct { //接收前端传过来的参数
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	//          接收前端参数，校验前端有没有传入字段，类型 格式，并将它映射到你的结构体里
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用 Service 层上传文档
	doc, err := h.documentService.UploadDocument(userID.(uint), req.Title, "", req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文档上传失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "文档上传成功",
		"data":    doc,
	})
}

// GetUserDocuments 获取当前用户的所有文档
func (h *DocumentHandler) GetUserDocuments(c *gin.Context) {
	userID, _ := c.Get("user_id")
	docs, err := h.documentService.GetUserDocuments(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文档列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": docs,
	})
}

// GetDocumentDetail 获取文档详情（带权限校验）
func (h *DocumentHandler) GetDocumentDetail(c *gin.Context) {
	userID, _ := c.Get("user_id")
	idStr := c.Param("id") //获取url上面的id
	//        把字符串转换为无符号整数       10进制    32位
	id, _ := strconv.ParseUint(idStr, 10, 32)

	doc, err := h.documentService.GetDocumentDetail(uint(id), userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文档不存在或无权限访问"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": doc,
	})
}

// DeleteDocument 删除文档（带权限校验）
func (h *DocumentHandler) DeleteDocument(c *gin.Context) {
	userID, _ := c.Get("user_id")
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	if err := h.documentService.DeleteDocument(uint(id), userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文档失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "文档删除成功",
	})
}
