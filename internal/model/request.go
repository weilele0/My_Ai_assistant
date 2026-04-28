package model

// 定义前端调用后端接口时，必须传递什么参数、格式是什么、哪些不能为空。
type RegisterRequest struct { //用户注册
	Username string `json:"username" binding:"required"` //不能为空
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type GenerateRequest struct { //ai生成内容
	Topic string `json:"topic" binding:"required"`
}

// RAGRequest RAG 生成请求
type RAGRequest struct {
	Question string `json:"question" binding:"required"`
	TopK     int    `json:"top_k" binding:"omitempty"` // 可选，默认返回 5 条
}

// RAGResponse RAG 生成响应
type RAGResponse struct {
	Answer     string   `json:"answer"`
	References []string `json:"references"` // 返回引用的文档片段
}
