# My AI Assistant

一个基于 **Go + Gin** 构建的全栈 AI 写作助手，集成 DeepSeek 文章生成、RAG 智能问答、文档管理与异步任务队列，前端采用原生 JS + 北欧极简风格设计。

---

## ✨ 核心功能

| 功能模块 | 说明 |
|----------|------|
| **AI 文章生成** | 输入主题，DeepSeek 大模型同步/异步生成高质量长文 |
| **RAG 智能问答** | 上传文档后，基于 Chroma 向量库 + 通义千问 Embedding 进行精准问答 |
| **文档管理** | 上传 / 查看 / 删除文档，文档自动向量化入库 |
| **异步任务队列** | 长文生成通过 Redis Stream 异步处理，实时轮询进度 |
| **额度管理** | 多套餐体系（免费 / 基础 / 专业 / 无限制），每日 AI & RAG 调用次数控制 |
| **管理后台** | 用户管理、套餐分配、系统统计、今日额度监控 |
| **用户系统** | JWT 鉴权，bcrypt 密码加密，管理员权限控制 |

---

## 🏗️ 技术栈

### 后端
- **语言 / 框架**：Go 1.26 + Gin
- **ORM**：GORM（自动迁移）
- **数据库**：MySQL（用户、文档、任务、订阅等）
- **缓存 / 队列**：Redis（会话 + Redis Stream 异步任务）
- **向量数据库**：Chroma（文档向量存储与检索）
- **AI 接口**：DeepSeek API（文章生成）、通义千问 API（Embedding）
- **鉴权**：JWT（golang-jwt/jwt v5）
- **配置管理**：godotenv（`.env` 文件）

### 前端
- **技术**：原生 HTML + Vanilla JS（无框架）
- **设计风格**：北欧极简 / 温暖陶土色系
- **布局**：SPA 单页应用，Gin 提供静态文件服务

---

## 📁 项目结构

```
My_AI_Assistant/
├── cmd/
│   └── main.go               # 程序入口，初始化所有组件
├── internal/
│   ├── config/               # 配置加载（读取 .env）
│   ├── db/                   # 数据库连接（MySQL / Redis / Chroma）
│   ├── handler/              # HTTP 控制器
│   │   ├── user_handler.go
│   │   ├── ai_handler.go
│   │   ├── document_handler.go
│   │   ├── task_handler.go
│   │   ├── rag_handler.go
│   │   ├── quota_handler.go
│   │   └── admin_handler.go
│   ├── middleware/            # 中间件（JWT / 额度校验 / 管理员权限）
│   ├── model/                # GORM 数据模型
│   ├── repository/           # 数据访问层
│   ├── router/               # 路由注册
│   └── service/              # 业务逻辑层
├── pkg/
│   └── jwt/                  # JWT 工具
├── frontend/
│   ├── index.html            # 单页入口
│   ├── css/style.css         # 全局样式
│   └── js/
│       ├── api.js            # 统一请求封装
│       ├── app.js            # 路由、Toast、模态框等全局函数
│       └── pages/            # 各页面模块（chat / documents / tasks / quota / admin / auth）
├── go.mod
├── go.sum
└── .env                      # 环境变量（需要自己创建）
```

---

## 🚀 快速开始

### 前置依赖

- Go 1.21+
- MySQL 8.0+
- Redis 6.0+
- [Chroma](https://www.trychroma.com/) 向量数据库（本地或 Docker 启动）
- DeepSeek API Key
- 通义千问 API Key

### 1. 克隆项目

```bash
git clone https://github.com/your-username/My_AI_Assistant.git
cd My_AI_Assistant
```

### 2. 配置环境变量

复制并修改 `.env` 文件：

```bash
cp .env.example .env
```

`.env` 内容示例：

```env
# AI 接口
DEEPSEEK_API_KEY=sk-xxxxxxxxxxxxxxxxxxxx
QWEN_API_KEY=sk-xxxxxxxxxxxxxxxxxxxx

# MySQL
MYSQL_DSN=root:password@tcp(127.0.0.1:3306)/ai_assistant?charset=utf8mb4&parseTime=True&loc=Local

# Redis
REDIS_ADDR=127.0.0.1:6379
REDIS_PASSWORD=
REDIS_DB=0

# Chroma 向量数据库
CHROMA_URL=http://localhost:8000

# JWT
JWT_SECRET=your-secret-key-here
```

### 3. 启动 Chroma（Docker）

```bash
docker run -d -p 8000:8000 chromadb/chroma
```

### 4. 安装依赖并启动

```bash
go mod tidy
go run cmd/main.go
```

> 首次启动会自动迁移数据库表结构，并初始化默认套餐数据。

### 5. 访问应用

打开浏览器访问：[http://localhost:8080](http://localhost:8080)

---

## 📄 项目文档

详细的技术方案和 API 接口说明请查阅 docs 目录：

- 📘 [技术方案文档](./docs/技术方案文档.md)
- 📗 [API 接口文档](./docs/API接口文档.md)

---

## 💰 套餐体系

| 套餐 | 每日 AI 生成 | 每日 RAG 问答 |
|------|-------------|--------------|
| 免费版（free） | 5 次 | 3 次 |
| 基础版（basic） | 30 次 | 20 次 |
| 专业版（pro） | 100 次 | 60 次 |
| 无限制（unlimit） | 不限 | 不限 |

> 管理员可通过后台为用户设置套餐，支持设置到期时间。

---

## ⚙️ 架构说明

### 异步任务流程
```
用户提交主题
    ↓
创建 AITask 记录（pending）
    ↓
任务 ID 写入 Redis Stream
    ↓
Worker Goroutine 消费 Stream
    ↓
调用 DeepSeek API 生成文章
    ↓
更新任务状态（completed / failed）
    ↓
前端轮询 /tasks/:id 获取结果
```

### RAG 流程
```
上传文档
    ↓
文档内容 → 通义千问 Embedding API → 向量
    ↓
向量存入 Chroma 集合 (documents_v2)
    ↓
用户提问 → 问题向量化
    ↓
Chroma 相似度检索 → 返回相关文档片段
    ↓
文档片段 + 用户问题 → DeepSeek → 最终回答
    ↓
保存到 RAG 历史记录
```

---

## 📝 数据库模型

| 表名 | 说明 |
|------|------|
| `users` | 用户信息（用户名、密码、是否管理员） |
| `documents` | 上传的文档（标题、内容、向量 ID） |
| `ai_tasks` | 异步任务（状态、主题、结果） |
| `rag_histories` | RAG 问答历史 |
| `subscriptions` | 套餐定义（每日 AI/RAG 上限） |
| `user_subscriptions` | 用户订阅（绑定套餐、到期时间） |
| `user_quotas` | 用户每日用量记录 |

---

## 🔒 安全说明

- 密码使用 **bcrypt** 加密存储，不可逆
- 所有业务接口均通过 **JWT 中间件**鉴权
- 管理员接口额外通过 **AdminAuth 中间件**验证
- 额度通过**中间件前置校验**，不足时直接拒绝请求

---

## 📄 License

MIT License
