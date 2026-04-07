package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config struct {
	ServerAddr string
	DBType     string // "mysql" or "sqlite"
	DBPath     string // sqlite file path or mysql connection string
}

func Load() *Config {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Println("未找到 .env 文件，使用系统环境变量")
	}

	cfg := &Config{
		ServerAddr: getEnv("SERVER_ADDR", ":8080"),
		DBType:     getEnv("DB_TYPE", "sqlite"),
		DBPath:     getEnv("DB_PATH", "policy_radar.db"),
	}
	return cfg
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func InitDB(cfg *Config) *gorm.DB {
	var db *gorm.DB
	var err error

	if cfg.DBType == "mysql" {
		// 解析连接字符串获取数据库名
		dbName := extractDatabaseName(cfg.DBPath)
		if dbName == "" {
			log.Fatalf("无法从连接字符串解析数据库名")
		}

		// 构建不带数据库名的连接字符串
		baseDSN := strings.Replace(cfg.DBPath, "/"+dbName, "/", 1)

		// 先连接到 MySQL 服务器（不指定数据库）
		baseDB, err := gorm.Open(mysql.Open(baseDSN), &gorm.Config{})
		if err != nil {
			log.Fatalf("MySQL 服务器连接失败: %v", err)
		}

		// 创建数据库（如果不存在）
		createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName)
		if err = baseDB.Exec(createDBSQL).Error; err != nil {
			log.Fatalf("数据库创建失败: %v", err)
		}
		log.Printf("数据库 '%s' 已准备好", dbName)

		// 关闭基础连接
		sqlDB, _ := baseDB.DB()
		sqlDB.Close()

		// 连接到目标数据库
		db, err = gorm.Open(mysql.Open(cfg.DBPath), &gorm.Config{})
	} else {
		db, err = gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	}

	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	return db
}

// extractDatabaseName 从 MySQL 连接字符串中提取数据库名
// DSN 格式: user:password@tcp(host:port)/dbname?params
func extractDatabaseName(dsn string) string {
	// 找到 / 后面的部分
	idx := strings.Index(dsn, "/")
	if idx == -1 {
		return ""
	}

	// 获取 / 后面的字符串
	rest := dsn[idx+1:]

	// 去掉问号后面的参数部分
	if qIdx := strings.Index(rest, "?"); qIdx != -1 {
		rest = rest[:qIdx]
	}

	return rest
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

type Pagination struct {
	Page     int    `form:"page" json:"page"`
PageSize  int    `form:"pageSize" json:"pageSize"`
	Keyword  string `form:"keyword" json:"keyword"`
Status    string `form:"status" json:"status"`
CheckStatus int   `form:"checkStatus" json:"checkStatus"`
	CategoryID  uint   `form:"categoryId" json:"categoryId"`
}

func (p *Pagination) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	return (p.Page - 1) * p.PageSize
}

func (p *Pagination) GetLimit() int {
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	return p.PageSize
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(data interface{}) Response {
	return Response{Code: 0, Message: "success", Data: data}
}

func Error(code int, message string) Response {
	return Response{Code: code, Message: message}
}

type PageResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    *PageData   `json:"data"`
}

type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
PageSize  int         `json:"pageSize"`
}