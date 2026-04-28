package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminAuth 管理员权限校验中间件
// 必须在 JWTAuth() 之后使用，依赖 JWTAuth 向 Context 写入的 is_admin 字段
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("is_admin")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
			c.Abort()
			return
		}

		// 类型断言，判断是否为管理员
		if admin, ok := isAdmin.(bool); !ok || !admin {
			c.JSON(http.StatusForbidden, gin.H{
				"code":  403,
				"error": "权限不足，仅管理员可访问",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
