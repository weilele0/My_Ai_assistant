# My_AI_Assistant API 接口文档

> 文档版本：v1.0.0
> 更新日期：2026-04-28
> Base URL：`http://localhost:8080`
> 所有接口返回格式均为 JSON

---

## 一、通用规范

### 1.1 认证方式

除公开接口外，所有接口均需携带 JWT Token：

```
Authorization: Bearer <token>
```

**获取 Token**：通过 `POST /api/v1/user/login` 登录后获取。

### 1.2 通用响应格式

**成功响应**：
```json
{
  "code": 200,
  "message": "操作成功描述",
  "data": { ... }
}
```

**错误响应**：
```json
{
  "code": 400,
  "error": "错误描述"
}
```

或
```json
{
  "code": 429,
  "error": "今日 AI 生成次数已达上限（5 次），请升级套餐"
}
```

### 1.3 HTTP 状态码

| 状态码 | 含义 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 / 业务校验失败 |
| 401 | 未认证（缺少 Token / Token 无效） |
| 403 | 权限不足（非管理员访问管理接口） |
| 404 | 资源不存在 |
| 429 | 额度不足 |
| 500 | 服务器内部错误 |

### 1.4 接口分组

| 前缀 | 认证要求 | 说明 |
|------|----------|------|
| `/api/v1/user/*` | 无 | 公开接口 |
| `/api/v1/ai/*` | JWT | AI 对话与 RAG |
| `/api/v1/documents/*` | JWT | 文档管理 |
| `/api/v1/tasks/*` | JWT | 异步任务 |
| `/api/v1/quota/*` | JWT | 额度查询 |
| `/api/v1/admin/*` | JWT + 管理员 | 后台管理 |

---

## 二、用户接口（公开）

### 2.1 用户注册

**POST** `/api/v1/user/register`

注册新用户。注册成功后需重新登录获取 Token。

**请求体**：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | ✅ | 用户名，2-50字符，唯一 |
| email | string | ✅ | 邮箱地址，格式校验 |
| password | string | ✅ | 密码，最少6字符 |

```json
{
  "username": "alice",
  "email": "alice@example.com",
  "password": "123456"
}
```

**成功响应** `(200)`：
```json
{
  "code": 200,
  "message": "注册成功"
}
```

**失败响应** `(400)`：
```json
{
  "error": "用户名已存在"
}
```

---

### 2.2 用户登录

**POST** `/api/v1/user/login`

登录并获取 JWT Token。

**请求体**：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | ✅ | 用户名 |
| password | string | ✅ | 密码 |

```json
{
  "username": "alice",
  "password": "123456"
}
```

**成功响应** `(200)`：
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "id": 1,
    "username": "alice",
    "is_admin": false
  }
}
```

**失败响应** `(401)`：
```json
{
  "error": "用户名或密码错误"
}
```

> ⚠️ Token 有效期为 **24 小时**，过期后需重新登录。

---

## 三、AI 接口（需认证）

### 3.1 同步 AI 生成

**POST** `/api/v1/ai/generate`

> ⚠️ 此接口会**消耗每日 AI 额度**（由 QuotaCheckAI 中间件自动扣减）

根据主题生成一篇约 800 字的高质量中文文章。

**请求头**：
```
Authorization: Bearer <token>
```

**请求体**：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| topic | string | ✅ | 文章主题 |

```json
{
  "topic": "人工智能在教育领域的应用与挑战"
}
```

**成功响应** `(200)`：
```json
{
  "code": 200,
  "message": "生成成功",
  "data": {
    "content": "引言：随着人工智能技术的快速发展，其在教育领域的应用日益广泛...\n\n核心内容：...\n\n总结：..."
  }
}
```

**失败响应** — 额度不足 `(429)`：
```json
{
  "code": 429,
  "error": "今日 AI 生成次数已达上限（5 次），请升级套餐"
}
```

**失败响应** — AI 服务错误 `(500)`：
```json
{
  "error": "AI 生成失败"
}
```

---

### 3.2 RAG 智能问答

**POST** `/api/v1/ai/rag-generate`

> ⚠️ 此接口会**消耗每日 RAG 额度**（由 QuotaCheckRAG 中间件自动扣减）

基于已上传文档进行智能问答，系统将检索相关文档片段并结合生成答案。

**前置条件**：需先通过 `POST /api/v1/documents` 上传至少一份文档。

**请求头**：
```
Authorization: Bearer <token>
```

**请求体**：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| question | string | ✅ | 用户提问 |
| top_k | int | ❌ | 检索的文档片段数量，默认 5 条 |

```json
{
  "question": "文档中提到的算法原理是什么？",
  "top_k": 5
}
```

**成功响应** `(200)`：
```json
{
  "code": 200,
  "data": {
    "answer": "根据参考资料，文档中提到的算法原理主要包括以下几个方面：\n\n【1】首先，...",
    "references": [
      "文档片段1的原文内容...",
      "文档片段2的原文内容...",
      "文档片段3的原文内容..."
    ]
  }
}
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| answer | string | AI 基于参考资料生成的答案，中文表述 |
| references | string[] | 检索到的相关文档片段原文，最多 top_k 条 |

**失败响应** — 无相关文档 `(200)`：
```json
{
  "code": 200,
  "data": {
    "answer": "根据现有参考资料无法回答该问题...",
    "references": []
  }
}
```

**失败响应** — 额度不足 `(429)`：
```json
{
  "code": 429,
  "error": "今日 RAG 问答次数已达上限（3 次），请升级套餐"
}
```

---

### 3.3 获取 RAG 历史列表

**GET** `/api/v1/ai/rag-history`

获取当前用户的 RAG 问答历史记录列表（按时间倒序）。

**请求头**：
```
Authorization: Bearer <token>
```

**成功响应** `(200)`：
```json
{
  "code": 200,
  "data": [
    {
      "id": 10,
      "user_id": 1,
      "question": "文档中提到的算法原理是什么？",
      "answer": "根据参考资料，...",
      "references": ["片段1", "片段2"],
      "created_at": "2026-04-28T14:30:00+08:00"
    },
    {
      "id": 9,
      "user_id": 1,
      "question": "...",
      "answer": "...",
      "references": [],
      "created_at": "2026-04-28T10:15:00+08:00"
    }
  ]
}
```

> ℹ️ 每位用户最多保留 **20 条** RAG 历史记录，超出后自动删除最旧的记录。

---

### 3.4 获取 RAG 历史详情

**GET** `/api/v1/ai/rag-history/:id`

获取单条 RAG 历史记录的完整内容。

**请求头**：
```
Authorization: Bearer <token>
```

**路径参数**：

| 参数 | 类型 | 说明 |
|------|------|------|
| :id | uint | 历史记录 ID |

**成功响应** `(200)`：
```json
{
  "code": 200,
  "data": {
    "id": 10,
    "user_id": 1,
    "question": "文档中提到的算法原理是什么？",
    "answer": "根据参考资料，文档中提到的算法原理主要包括以下几个方面：\n\n【1】...",
    "references": ["片段1原文内容...", "片段2原文内容..."],
    "created_at": "2026-04-28T14:30:00+08:00"
  }
}
```

**失败响应** `(404)`：
```json
{
  "error": "历史记录不存在或无权限"
}
```

---

### 3.5 删除单条 RAG 历史

**DELETE** `/api/v1/ai/rag-history/:id`

**请求头**：
```
Authorization: Bearer <token>
```

**路径参数**：

| 参数 | 类型 | 说明 |
|------|------|------|
| :id | uint | 历史记录 ID |

**成功响应** `(200)`：
```json
{
  "code": 200,
  "message": "删除成功"
}
```

---

### 3.6 清空所有 RAG 历史

**DELETE** `/api/v1/ai/rag-history`

清空当前用户的所有 RAG 历史记录。

**请求头**：
```
Authorization: Bearer <token>
```

**成功响应** `(200)`：
```json
{
  "code": 200,
  "message": "已清空所有历史记录"
}
```

---

## 四、文档接口（需认证）

### 4.1 上传文档

**POST** `/api/v1/documents`

上传纯文本内容，自动向量化并存入 Chroma 向量数据库。

> ℹ️ 向量化过程在 Goroutine 中异步执行，接口立即返回。上传后需等待数秒再进行 RAG 问答。

**请求头**：
```
Authorization: Bearer <token>
```

**请求体**：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | ✅ | 文档标题 |
| content | string | ✅ | 文档正文内容（纯文本） |

```json
{
  "title": "深度学习算法原理",
  "content": "一、神经网络基础\n神经网络是一种模拟人脑神经元运作的计算模型..."
}
```

**成功响应** `(200)`：
```json
{
  "code": 200,
  "message": "文档上传成功",
  "data": {
    "id": 1,
    "user_id": 1,
    "title": "深度学习算法原理",
    "file_path": "",
    "content": "一、神经网络基础...",
    "created_at": "2026-04-28T14:00:00+08:00",
    "updated_at": "2026-04-28T14:00:00+08:00"
  }
}
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 文档 ID，后续用于 RAG 检索 |
| file_path | string | 文件路径（当前版本为空） |
| content | string | 原始文本内容 |

---

### 4.2 获取文档列表

**GET** `/api/v1/documents`

获取当前用户上传的所有文档（按上传时间倒序）。

**请求头**：
```
Authorization: Bearer <token>
```

**成功响应** `(200)`：
```json
{
  "code": 200,
  "data": [
    {
      "id": 2,
      "user_id": 1,
      "title": "Go 语言实战笔记",
      "file_path": "",
      "content": "...",
      "created_at": "2026-04-28T13:00:00+08:00",
      "updated_at": "2026-04-28T13:00:00+08:00"
    },
    {
      "id": 1,
      "user_id": 1,
      "title": "深度学习算法原理",
      "file_path": "",
      "content": "...",
      "created_at": "2026-04-28T14:00:00+08:00",
      "updated_at": "2026-04-28T14:00:00+08:00"
    }
  ]
}
```

---

### 4.3 获取文档详情

**GET** `/api/v1/documents/:id`

获取指定文档的完整内容。

**请求头**：
```
Authorization: Bearer <token>
```

**路径参数**：

| 参数 | 类型 | 说明 |
|------|------|------|
| :id | uint | 文档 ID |

**成功响应** `(200)`：
```json
{
  "code": 200,
  "data": {
    "id": 1,
    "user_id": 1,
    "title": "深度学习算法原理",
    "file_path": "",
    "content": "一、神经网络基础\n神经网络是一种模拟人脑神经元运作的计算模型...",
    "created_at": "2026-04-28T14:00:00+08:00",
    "updated_at": "2026-04-28T14:00:00+08:00"
  }
}
```

**失败响应** `(404)`：
```json
{
  "error": "文档不存在或无权限访问"
}
```

---

### 4.4 删除文档

**DELETE** `/api/v1/documents/:id`

删除指定文档（同时删除 MySQL 记录和 Chroma 向量，**注意**：当前版本仅删除 MySQL 记录，Chroma 中的向量需手动清理）。

**请求头**：
```
Authorization: Bearer <token>
```

**路径参数**：

| 参数 | 类型 | 说明 |
|------|------|------|
| :id | uint | 文档 ID |

**成功响应** `(200)`：
```json
{
  "code": 200,
  "message": "文档删除成功"
}
```

---

## 五、任务接口（需认证）

### 5.1 提交异步生成任务

**POST** `/api/v1/tasks/generate`

提交长文生成任务，立即返回任务 ID，生成过程在后台异步执行。

> ℹ️ 此接口提交任务时已扣除额度，任务无论成功失败均消耗一次。

**请求头**：
```
Authorization: Bearer <token>
```

**请求体**：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| topic | string | ✅ | 文章主题 |

```json
{
  "topic": "区块链技术在供应链管理中的应用"
}
```

**成功响应** `(200)`：
```json
{
  "code": 200,
  "message": "任务已提交，正在后台处理",
  "data": {
    "task_id": 5
  }
}
```

---

### 5.2 获取任务列表

**GET** `/api/v1/tasks`

获取当前用户的所有异步任务（按创建时间倒序）。

**请求头**：
```
Authorization: Bearer <token>
```

**成功响应** `(200)`：
```json
{
  "code": 200,
  "data": [
    {
      "id": 5,
      "user_id": 1,
      "task_type": "generate",
      "input_text": "区块链技术在供应链管理中的应用",
      "output_text": "引言：...",
      "status": "completed",
      "error_msg": "",
      "created_at": "2026-04-28T15:00:00+08:00",
      "completed_at": "2026-04-28T15:00:35+08:00"
    },
    {
      "id": 4,
      "user_id": 1,
      "task_type": "generate",
      "input_text": "量子计算的未来",
      "output_text": "",
      "status": "processing",
      "error_msg": "",
      "created_at": "2026-04-28T14:50:00+08:00",
      "completed_at": null
    }
  ]
}
```

**任务状态说明**：

| status | 说明 |
|--------|------|
| `pending` | 任务已创建，等待 Worker 消费 |
| `processing` | Worker 正在处理中 |
| `completed` | 生成完成，output_text 为生成结果 |
| `failed` | 生成失败，error_msg 为错误原因 |

---

### 5.3 查询任务状态

**GET** `/api/v1/tasks/:id`

查询指定任务的状态和结果。

**请求头**：
```
Authorization: Bearer <token>
```

**路径参数**：

| 参数 | 类型 | 说明 |
|------|------|------|
| :id | uint | 任务 ID |

**成功响应 — 处理中** `(200)`：
```json
{
  "code": 200,
  "data": {
    "id": 5,
    "user_id": 1,
    "task_type": "generate",
    "input_text": "区块链技术在供应链管理中的应用",
    "output_text": "",
    "status": "processing",
    "error_msg": "",
    "created_at": "2026-04-28T15:00:00+08:00",
    "completed_at": null
  }
}
```

**成功响应 — 已完成** `(200)`：
```json
{
  "code": 200,
  "data": {
    "id": 5,
    "user_id": 1,
    "task_type": "generate",
    "input_text": "区块链技术在供应链管理中的应用",
    "output_text": "引言：随着区块链技术的不断成熟，其在供应链管理领域的应用前景日益受到关注...\n\n...",
    "status": "completed",
    "error_msg": "",
    "created_at": "2026-04-28T15:00:00+08:00",
    "completed_at": "2026-04-28T15:00:35+08:00"
  }
}
```

**失败响应 — 任务失败** `(200)`：
```json
{
  "code": 200,
  "data": {
    "id": 5,
    "status": "failed",
    "error_msg": "AI 服务调用超时",
    "output_text": ""
  }
}
```

---

## 六、额度接口（需认证）

### 6.1 查看我的额度

**GET** `/api/v1/quota/me`

查看当前用户的套餐信息和今日用量。

**请求头**：
```
Authorization: Bearer <token>
```

**成功响应** `(200)`：
```json
{
  "code": 200,
  "data": {
    "plan": "free",
    "daily_ai_quota": 5,
    "daily_rag_quota": 3,
    "used_ai": 2,
    "used_rag": 1,
    "remain_ai": 3,
    "remain_rag": 2,
    "expired_at": null
  }
}
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| plan | string | 当前套餐标识（free/basic/pro/unlimit） |
| daily_ai_quota | int | 每日 AI 生成上限（-1=无限） |
| daily_rag_quota | int | 每日 RAG 问答上限（-1=无限） |
| used_ai | int | 今日已用 AI 生成次数 |
| used_rag | int | 今日已用 RAG 问答次数 |
| remain_ai | int | 今日剩余 AI 次数（-1=无限） |
| remain_rag | int | 今日剩余 RAG 次数（-1=无限） |
| expired_at | string/null | 套餐到期时间（null=永久有效） |

---

## 七、管理后台接口（需认证 + 管理员权限）

> ⚠️ 以下接口需要 `is_admin=true`，普通用户访问将返回 `403` 权限不足。

### 7.1 用户列表

**GET** `/api/v1/admin/users`

获取所有注册用户列表。

**请求头**：
```
Authorization: Bearer <token>  （需 is_admin=true）
```

**成功响应** `(200)`：
```json
{
  "code": 200,
  "data": [
    {
      "id": 1,
      "username": "alice",
      "email": "alice@example.com",
      "is_admin": true,
      "created_at": "2026-04-20T10:00:00+08:00"
    },
    {
      "id": 2,
      "username": "bob",
      "email": "bob@example.com",
      "is_admin": false,
      "created_at": "2026-04-25T14:30:00+08:00"
    }
  ]
}
```

---

### 7.2 用户详情与额度

**GET** `/api/v1/admin/users/:id`

查看指定用户的详细信息和今日额度使用情况。

**请求头**：
```
Authorization: Bearer <token>  （需 is_admin=true）
```

**路径参数**：

| 参数 | 类型 | 说明 |
|------|------|------|
| :id | uint | 用户 ID |

**成功响应** `(200)`：
```json
{
  "code": 200,
  "data": {
    "user": {
      "id": 2,
      "username": "bob",
      "email": "bob@example.com",
      "is_admin": false,
      "created_at": "2026-04-25T14:30:00+08:00"
    },
    "quota": {
      "plan": "basic",
      "daily_ai_quota": 30,
      "daily_rag_quota": 20,
      "used_ai": 5,
      "used_rag": 3,
      "remain_ai": 25,
      "remain_rag": 17,
      "expired_at": null
    }
  }
}
```

---

### 7.3 设置用户套餐

**POST** `/api/v1/admin/users/:id/plan`

为指定用户设置或变更套餐。

**请求头**：
```
Authorization: Bearer <token>  （需 is_admin=true）
```

**路径参数**：

| 参数 | 类型 | 说明 |
|------|------|------|
| :id | uint | 用户 ID |

**请求体**：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| plan | string | ✅ | 套餐标识：free / basic / pro / unlimit |
| expired_at | string | ❌ | 到期时间（ISO 8601 格式），设为 null 表示永久有效 |

```json
{
  "plan": "pro",
  "expired_at": "2026-12-31T23:59:59Z"
}
```

**成功响应** `(200)`：
```json
{
  "code": 200,
  "message": "套餐设置成功"
}
```

**失败响应** `(400)`：
```json
{
  "error": "无效的套餐类型: invalid_plan"
}
```

---

### 7.4 系统统计

**GET** `/api/v1/admin/stats`

获取系统运营统计数据。

**请求头**：
```
Authorization: Bearer <token>  （需 is_admin=true）
```

**成功响应** `(200)`：
```json
{
  "code": 200,
  "data": {
    "total_users": 42,
    "today_active_users": 15,
    "today_ai_calls": 87,
    "today_rag_calls": 53,
    "plan_distribution": {
      "free": 30,
      "basic": 8,
      "pro": 3,
      "unlimit": 1
    },
    "today_quota_details": [
      {
        "id": 100,
        "user_id": 1,
        "date": "2026-04-28",
        "used_ai": 12,
        "used_rag": 5
      }
    ]
  }
}
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| total_users | int | 注册用户总数 |
| today_active_users | int | 今日有额度使用记录的用户数 |
| today_ai_calls | int | 今日所有用户 AI 生成总调用次数 |
| today_rag_calls | int | 今日所有用户 RAG 问答总调用次数 |
| plan_distribution | object | 各套餐用户数量分布 |
| today_quota_details | array | 今日所有用户额度明细（分页可扩展） |

---

### 7.5 今日额度明细

**GET** `/api/v1/admin/quotas/today`

查看今日所有用户的额度使用情况明细。

**请求头**：
```
Authorization: Bearer <token>  （需 is_admin=true）
```

**成功响应** `(200)`：
```json
{
  "code": 200,
  "data": [
    {
      "id": 101,
      "user_id": 1,
      "date": "2026-04-28",
      "used_ai": 12,
      "used_rag": 5
    },
    {
      "id": 102,
      "user_id": 2,
      "date": "2026-04-28",
      "used_ai": 5,
      "used_rag": 3
    }
  ]
}
```

> ℹ️ 返回结果按 `used_ai + used_rag` 倒序排列，活跃用户排在前面。

---

### 7.6 订阅列表

**GET** `/api/v1/admin/subscriptions`

查看所有用户的当前订阅情况。

**请求头**：
```
Authorization: Bearer <token>  （需 is_admin=true）
```

**成功响应** `(200)`：
```json
{
  "code": 200,
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "plan": "unlimit",
      "expired_at": null,
      "created_at": "2026-04-20T10:00:00+08:00",
      "updated_at": "2026-04-25T09:00:00+08:00"
    },
    {
      "id": 2,
      "user_id": 2,
      "plan": "basic",
      "expired_at": "2026-12-31T23:59:59Z",
      "created_at": "2026-04-25T14:30:00+08:00",
      "updated_at": "2026-04-25T14:30:00+08:00"
    }
  ]
}
```

---

## 八、完整接口索引

| # | 方法 | 路径 | 认证 | 额度消耗 | 说明 |
|---|------|------|------|----------|------|
| 1 | POST | `/api/v1/user/register` | ❌ | — | 用户注册 |
| 2 | POST | `/api/v1/user/login` | ❌ | — | 用户登录 |
| 3 | POST | `/api/v1/ai/generate` | ✅ | AI | 同步 AI 生成文章 |
| 4 | POST | `/api/v1/ai/rag-generate` | ✅ | RAG | RAG 智能问答 |
| 5 | GET | `/api/v1/ai/rag-history` | ✅ | — | 获取 RAG 历史列表 |
| 6 | GET | `/api/v1/ai/rag-history/:id` | ✅ | — | 获取 RAG 历史详情 |
| 7 | DELETE | `/api/v1/ai/rag-history/:id` | ✅ | — | 删除单条 RAG 历史 |
| 8 | DELETE | `/api/v1/ai/rag-history` | ✅ | — | 清空所有 RAG 历史 |
| 9 | POST | `/api/v1/documents` | ✅ | — | 上传文档 |
| 10 | GET | `/api/v1/documents` | ✅ | — | 获取文档列表 |
| 11 | GET | `/api/v1/documents/:id` | ✅ | — | 获取文档详情 |
| 12 | DELETE | `/api/v1/documents/:id` | ✅ | — | 删除文档 |
| 13 | POST | `/api/v1/tasks/generate` | ✅ | AI | 提交异步生成任务 |
| 14 | GET | `/api/v1/tasks` | ✅ | — | 获取任务列表 |
| 15 | GET | `/api/v1/tasks/:id` | ✅ | — | 查询任务状态 |
| 16 | GET | `/api/v1/quota/me` | ✅ | — | 查看我的额度 |
| 17 | GET | `/api/v1/admin/users` | ✅+Admin | — | 用户列表 |
| 18 | GET | `/api/v1/admin/users/:id` | ✅+Admin | — | 用户详情与额度 |
| 19 | POST | `/api/v1/admin/users/:id/plan` | ✅+Admin | — | 设置用户套餐 |
| 20 | GET | `/api/v1/admin/stats` | ✅+Admin | — | 系统统计 |
| 21 | GET | `/api/v1/admin/quotas/today` | ✅+Admin | — | 今日额度明细 |
| 22 | GET | `/api/v1/admin/subscriptions` | ✅+Admin | — | 订阅列表 |

---

## 九、套餐额度一览

| 套餐 | plan 值 | 每日 AI 上限 | 每日 RAG 上限 |
|------|---------|-------------|--------------|
| 免费版 | `free` | 5 次 | 3 次 |
| 基础版 | `basic` | 30 次 | 20 次 |
| 专业版 | `pro` | 100 次 | 60 次 |
| 无限制 | `unlimit` | 不限 (-1) | 不限 (-1) |

> ℹ️ 额度按自然日统计，每日凌晨重置。过期的订阅自动降级为免费版。

---

*本文档随接口变更持续更新*
