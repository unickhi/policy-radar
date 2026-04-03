package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"policy-radar/internal/config"
	"policy-radar/internal/model"
	"policy-radar/internal/repository"
	"policy-radar/internal/service"
)

type RecommendHandler struct {
	svc *service.RecommendService
}

func NewRecommendHandler(svc *service.RecommendService) *RecommendHandler {
	return &RecommendHandler{svc: svc}
}

func (h *RecommendHandler) List(c *gin.Context) {
	// 获取查询参数
	keyword := c.Query("keyword")
	policyType := c.Query("policy_type")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	query := &repository.RecommendQuery{
		Keyword:    keyword,
		PolicyType: policyType,
		Page:       page,
		PageSize:   pageSize,
	}

	list, total, err := h.svc.List(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, config.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, config.PageResponse{
		Code:    0,
		Message: "success",
		Data: &config.PageData{
			List:    list,
			Total:   total,
			Page:    page,
			PageSize: pageSize,
		},
	})
}

func (h *RecommendHandler) Create(c *gin.Context) {
	var item model.PolicyRecommend
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, config.Error(400, "参数错误"))
		return
	}

	// 使用 Upsert，如果有 id 则更新，否则创建
	if err := h.svc.Upsert(&item); err != nil {
		c.JSON(http.StatusInternalServerError, config.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, config.Success(item))
}

func (h *RecommendHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, config.Error(400, "ID无效"))
		return
	}

	var item model.PolicyRecommend
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, config.Error(400, "参数错误"))
		return
	}

	item.ID = uint(id)
	if err := h.svc.Update(&item); err != nil {
		c.JSON(http.StatusInternalServerError, config.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, config.Success(item))
}

func (h *RecommendHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, config.Error(400, "ID无效"))
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, config.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, config.Success(gin.H{"message": "删除成功"}))
}