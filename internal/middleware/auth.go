package middleware

import (
	"My_AI_Assistant/pkg/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT 认证中间件
func JWTAuth() gin.HandlerFunc { //Gin 中间件构造函数
	return func(c *gin.Context) {
		// 1. 从请求头获取 Authorization
		//             获取HTTP请求头字段的值     HTTP 标准认证头字段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "请提供 Authorization 头"})
			c.Abort()
			return
		}

		// 2. 解析 Bearer Token
		//        把一段字符串按分隔符切开  传入待切割的字符串    分隔符   最大切成两个
		parts := strings.SplitN(authHeader, " ", 2)
		//       切成两段          第一段是Bearer
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization 格式错误，应为 Bearer <token>"})
			c.Abort() // 终止当前请求的所有后续处理，直接返回响应
			return
		}

		// 3. 解析并验证 Token
		//              解析前端的token     第一段是头部：bearer 第二段是token
		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token 无效或已过期"})
			c.Abort()
			return
		}

		// 4. 将用户信息存入 Context，供后续 Handler 使用
		//Gin 的请求上下文（Context）里存数据，让后面所有 handler 都能读到。
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("is_admin", claims.IsAdmin)

		// 5. 继续执行下一个 Handler
		//运行逻辑   请求-》 鉴权-》 业务处理-》 响应
		c.Next()
	}
}
