package handler

import (
	"My_AI_Assistant/internal/model"
	"My_AI_Assistant/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TaskHandler 任务控制器
type TaskHandler struct {
	taskService *service.TaskService
}

// NewTaskHandler 创建 TaskHandler
func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

// CreateAsyncGenerateTask 创建异步生成任务
func (h *TaskHandler) CreateAsyncGenerateTask(c *gin.Context) {
	userID, _ := c.Get("user_id") // 从 JWT 中间件获取

	var req model.GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建异步任务（立即返回任务ID）
	taskID, err := h.taskService.CreateAsyncTask(userID.(uint), "generate", req.Topic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建任务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "任务已提交，正在后台处理",
		"data": gin.H{
			"task_id": taskID,
		},
	})
}

// GetTaskStatus 查询任务状态
func (h *TaskHandler) GetTaskStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")
	taskIDStr := c.Param("id")
	taskID, _ := strconv.ParseUint(taskIDStr, 10, 32) //转为uint64类型

	task, err := h.taskService.GetTaskStatus(uint(taskID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在或无权限查看"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": task,
	})
}

// GetUserTasks 获取当前用户的所有任务列表
func (h *TaskHandler) GetUserTasks(c *gin.Context) {
	userID, _ := c.Get("user_id")

	tasks, err := h.taskService.GetUserTasks(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取任务列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": tasks,
	})
}
