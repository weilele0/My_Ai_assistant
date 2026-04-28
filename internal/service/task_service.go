package service

import (
	"My_AI_Assistant/internal/model"
	"My_AI_Assistant/internal/repository"
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type TaskService struct {
	aiService   *AIService                   //AI 大模型调用服务
	taskRepo    *repository.AITaskRepository //任务数据访问层
	redisClient *redis.Client                //Redis 客户端，用于操作 Stream 消息队列
}

// NewTaskService 创建 TaskService
func NewTaskService(aiService *AIService, taskRepo *repository.AITaskRepository, redisClient *redis.Client) *TaskService {
	return &TaskService{
		aiService:   aiService,
		taskRepo:    taskRepo,
		redisClient: redisClient,
	}
}

// CreateAsyncTask 创建异步任务（发布到 Redis Stream）
func (s *TaskService) CreateAsyncTask(userID uint, taskType, inputText string) (uint, error) {
	//初始化任务实体，状态设置为 pending
	task := &model.AITask{
		UserID:    userID,
		TaskType:  taskType,
		InputText: inputText,
		Status:    "pending",
	}

	if err := s.taskRepo.Create(task); err != nil {
		return 0, err
	}

	// 发布任务到 Redis Stream
	//Redis Stream 消息本身就是键值对结构必须用map类型
	taskData := map[string]interface{}{ //用于封装要发送到 Redis Stream 的任务数据
		"task_id":   task.ID,
		"user_id":   userID,
		"task_type": taskType,
		"input":     inputText,
	}
	//s.redisClient.XAdd(...)写入消息命令  上下文          go-redis 客户端封装的参数结构体，表示要发送的消息内容
	_, err := s.redisClient.XAdd(context.Background(), &redis.XAddArgs{
		Stream: "ai_task_stream", //发送到哪里
		Values: taskData,         //消息内容
	}).Result()

	if err != nil {
		return 0, fmt.Errorf("发布任务到 Redis Stream 失败: %v", err)
	}

	return task.ID, nil
}

// StartTaskWorker 启动后台 Worker 监听 Redis Stream
// 后台一直运行，监听消息队列，有任务就处理，没有就阻塞等待
func (s *TaskService) StartTaskWorker() {
	ctx := context.Background()      //上下文
	consumerGroup := "ai_task_group" //消费者组名
	consumerName := "worker-1"       //消费者名称

	// 创建消费者组（如果不存在） 给一个 Stream 创建一个消费者组 // 要操作的 Stream 队列名字          从哪里开始消费
	s.redisClient.XGroupCreateMkStream(ctx, "ai_task_stream", consumerGroup, "0").Err()

	fmt.Println(">>> Redis Stream Worker 已启动，监听 ai_task_stream")

	for { //           Redis 消费者组读取消息命令            封装所有读取消息需要的参数
		streams, err := s.redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    consumerGroup,                   //哪个组
			Consumer: consumerName,                    //谁
			Streams:  []string{"ai_task_stream", ">"}, // // 读哪个队列，读新消息
			Count:    1,                               // 一次读几条
			Block:    0,                               // 阻塞等待
		}).Result()

		if err != nil {
			fmt.Println(">>> Redis Stream 读取失败:", err)
			time.Sleep(1 * time.Second) //休息一秒
			continue                    //重试
		}

		for _, stream := range streams { //遍历消息队列，取出一个消息
			for _, message := range stream.Messages { //遍历这个消息盒子里具体的消息
				s.processStreamMessage(message) //处理消息
			}
		}
	}
}

// processStreamMessage 处理从 Stream 接收到的消息
func (s *TaskService) processStreamMessage(msg redis.XMessage) {
	taskIDStr := msg.Values["task_id"].(string) //从消息的键值对中取出 task_id
	taskID := uint(0)
	//把字符串类型的转换为数字     解析成整数     放入这个里面
	fmt.Sscanf(taskIDStr, "%d", &taskID)

	task, err := s.taskRepo.FindByIDForWorker(taskID)
	if err != nil {
		return
	}

	// 更新状态为 processing
	task.Status = "processing"
	s.taskRepo.Update(task)

	// 调用 AI 生成内容
	content, err := s.aiService.GenerateContent(task.InputText)

	if err != nil {
		task.Status = "failed"
		task.ErrorMsg = err.Error()
	} else {
		task.Status = "completed"
		task.OutputText = content
	}

	now := time.Now()
	task.CompletedAt = &now
	s.taskRepo.Update(task)

	// 确认消息已处理
	//               标记完成                            哪个队列               哪个组         消息id
	s.redisClient.XAck(context.Background(), "ai_task_stream", "ai_task_group", msg.ID)
}

// GetTaskStatus 查询任务状态（带用户权限校验）
func (s *TaskService) GetTaskStatus(taskID uint, userID uint) (*model.AITask, error) {
	return s.taskRepo.FindByID(taskID, userID)
}

// GetUserTasks 获取用户所有任务（最新在前）
func (s *TaskService) GetUserTasks(userID uint) ([]model.AITask, error) {
	return s.taskRepo.FindByUserID(userID)
}
