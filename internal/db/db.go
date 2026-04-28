package db

import (
	"My_AI_Assistant/internal/config"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitMySQL 初始化 MySQL 连接
func InitMySQL() *gorm.DB {
	cfg := config.LoadConfig()

	fmt.Print("正在连接 MySQL... ")
	db, err := gorm.Open(mysql.Open(cfg.MySQLDSN), &gorm.Config{})
	if err != nil {
		panic("❌ MySQL 连接失败: " + err.Error())
	}
	fmt.Println("✅ MySQL 连接成功")

	// 设置连接池
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db
}

// InitRedis 初始化 Redis 连接
func InitRedis() *redis.Client {
	cfg := config.LoadConfig()

	fmt.Print("正在连接 Redis... ")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	_, err := redisClient.Ping(redisClient.Context()).Result()
	if err != nil {
		panic("❌ Redis 连接失败: " + err.Error())
	}
	fmt.Println("✅ Redis 连接成功")

	return redisClient
}

// InitChroma 检查 Chroma 连接（严格模式）
func InitChroma() {
	cfg := config.LoadConfig()

	fmt.Print("正在检查 Chroma 连接... ")
	chromaOK := false

	for i := 0; i < 12; i++ {
		resp, err := http.Get(cfg.ChromaURL + "/api/v1/heartbeat")
		if err == nil && resp.StatusCode == 200 {
			chromaOK = true
			resp.Body.Close()
			break
		}
		if resp != nil {
			resp.Body.Close()
		}
		fmt.Print(".")
		time.Sleep(2 * time.Second)
	}

	if !chromaOK {
		panic("\n❌ Chroma 连接失败！请确认 Chroma 容器是否正常运行")
	}
	fmt.Println(" ✅ Chroma 连接成功")
}
