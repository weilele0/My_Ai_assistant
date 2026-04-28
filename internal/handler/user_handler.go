package handler

import (
	"My_AI_Assistant/internal/model"
	"My_AI_Assistant/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户控制器
type UserHandler struct {
	userService *service.UserService //给控制器（Handler）装上 “业务工具”
}

// NewUserHandler 创建 UserHandler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Register 处理用户注册请求
func (h *UserHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	//把前端传过来的 JSON 数据，自动解析绑定到结构体里，方便后面使用
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//调用服务层注册
	if err := h.userService.Register(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "注册成功",
	})
}

// Login 处理用户登录请求
func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.userService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data": gin.H{
			"token":   token,
			"id":      user.ID,
			"username": user.Username,
			"is_admin": user.IsAdmin,
		},
	})
}
