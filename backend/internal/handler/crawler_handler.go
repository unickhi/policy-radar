package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"policy-radar/internal/config"
	"policy-radar/internal/crawler"
	"policy-radar/internal/service"
)

type CrawlerHandler struct {
	svc *service.CrawlerService
}

func NewCrawlerHandler(svc *service.CrawlerService) *CrawlerHandler {
	return &CrawlerHandler{svc: svc}
}

func (h *CrawlerHandler) Execute(c *gin.Context) {
	var req struct {
		Script string `json:"script"`
		Query  string `json:"query"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, config.Error(400, "参数错误"))
		return
	}

	// 验证脚本安全性
	if err := crawler.ValidateScript(req.Script); err != nil {
		c.JSON(http.StatusBadRequest, config.Error(400, err.Error()))
		return
	}

	// 执行脚本
	result, log, err := crawler.ExecutePython(req.Script, req.Query)
	if err != nil {
		c.JSON(http.StatusOK, config.Success(gin.H{
			"success": false,
			"log":     log,
			"error":   err.Error(),
			"data":    []interface{}{},
		}))
		return
	}

	c.JSON(http.StatusOK, config.Success(gin.H{
		"success": true,
		"log":     log,
		"data":    result,
	}))
}

func (h *CrawlerHandler) Logs(c *gin.Context) {
	// TODO: 获取执行日志
	c.JSON(http.StatusOK, config.Success([]interface{}{}))
}

func (h *CrawlerHandler) Import(c *gin.Context) {
	var req struct {
		Data       []map[string]interface{} `json:"data"`
		TargetType string                   `json:"target_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, config.Error(400, "参数错误"))
		return
	}

	count, err := h.svc.ImportData(req.Data, req.TargetType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, config.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, config.Success(gin.H{"count": count}))
}