package jwt

import (
	"My_AI_Assistant/internal/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT 自定义声明
type Claims struct {
	UserID               uint   `json:"user_id"` //用户id
	Username             string `json:"username"`
	IsAdmin              bool   `json:"is_admin"` // 是否为管理员
	jwt.RegisteredClaims        // JWT 标准字段（过期时间、签发时间等）
	/*过期时间
	签发时间
	生效时间
	作用：控制 Token 多久过期*/
}

// GenerateToken 生成 JWT Token
func GenerateToken(userID uint, username string, isAdmin bool) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24小时过期//过期时间  过期后，这个 Token 就作废了，用户必须重新登录
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     //签发时间 作用：记录这个 Token 是什么时候生成的
			Issuer:    "My_AI_Assistant",                                  //签发人 作用：标识这个 Token 是谁发的
		},
	}
	// 打包，生成一个待签名的jwt对象 对称加密算法 保证token不被篡改  用户信息
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//           签名+加密             //调用这个方法  获取密钥
	return token.SignedString([]byte(config.LoadConfig().JWTSecret))
}

// ParseToken 解析并验证 Token
func ParseToken(tokenString string) (*Claims, error) {
	//jwt.token类型             token解析函数        前端传入的token    解析到结构体   回调函数 自动传入正在解析的token   jwt支持多种签名算法，返回值不同
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		//        密钥
		return []byte(config.LoadConfig().JWTSecret), nil
		//因为HS256这个算法底层只接受字节数组
		//解析到了token内部。直接赋值给 token.Claims 接口字段
	})
	//上面是解析token，并将解析后的token返回到结构体里
	if err != nil {
		return nil, err
	}
	//               断言 //这个接口是不是claims类型 是的话，把它还原成原来的结构体指针
	if claims, ok := token.Claims.(*Claims); ok && token.Valid { //Token 是否合法
		return claims, nil ////////////////断言是否成功  Token 是否合法
	}
	//token.claims这个库自带的接口里面有数据，但是系统不知到，只有判断一下才可以取出来
	return nil, errors.New("token validation failed")
}
