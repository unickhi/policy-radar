package config

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
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
		db, err = gorm.Open(mysql.Open(cfg.DBPath), &gorm.Config{})
	} else {
		db, err = gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	}

	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	return db
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