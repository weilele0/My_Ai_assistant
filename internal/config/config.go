package config

import (
	"os"

	"github.com/joho/godotenv"
)

// go项目配置加载代码，用来读取.env环境变量
type Config struct {
	DeepSeekAPIKey string //AI助手的密钥
	QwenAPIKey     string
	MySQLDSN       string //Mysql数据库连接的地址
	RedisAddr      string //redis地址
	RedisPassword  string //密码
	RedisDB        int    //用第几个库
	JWTSecret      string //登录验证的密钥
	ChromaURL      string
}

func LoadConfig() *Config {
	godotenv.Load() //加载。env文件
	return &Config{ //返回一个填好数据的config结构体
		DeepSeekAPIKey: os.Getenv("DEEPSEEK_API_KEY"), //读取后面的关键字，给前面赋值
		QwenAPIKey:     os.Getenv("QWEN_API_KEY"),
		MySQLDSN:       os.Getenv("MYSQL_DSN"),
		RedisAddr:      os.Getenv("REDIS_ADDR"),
		RedisPassword:  os.Getenv("REDIS_PASSWORD"),
		RedisDB:        0,
		JWTSecret:      os.Getenv("JWT_SECRET"),
		ChromaURL:      os.Getenv("CHROMA_URL"),
	}
}
