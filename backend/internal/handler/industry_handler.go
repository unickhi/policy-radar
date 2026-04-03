package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"policy-radar/internal/config"
	"policy-radar/internal/model"
	"policy-radar/internal/service"
)

type IndustryStandardHandler struct {
	svc *service.IndustryStandardService
}

func NewIndustryStandardHandler(svc *service.IndustryStandardService) *IndustryStandardHandler {
	return &IndustryStandardHandler{svc: svc}
}

func (h *IndustryStandardHandler) List(c *gin.Context) {
	var p config.Pagination
	if err := c.ShouldBindQuery(&p); err != nil {
		c.JSON(http.StatusBadRequest, config.Error(400, "参数错误"))
		return
	}

	list, total, err := h.svc.List(&p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, config.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, config.PageResponse{
		Code:    0,
		Message: "success",
		Data: &config.PageData{
			List:     list,
			Total:    total,
			Page:     p.Page,
PageSize:  p.PageSize,
		},
	})
}

func (h *IndustryStandardHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, config.Error(400, "ID无效"))
		return
	}

	item, err := h.svc.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, config.Error(404, "记录不存在"))
		return
	}

	c.JSON(http.StatusOK, config.Success(item))
}

func (h *IndustryStandardHandler) Create(c *gin.Context) {
	var item model.IndustryStandard
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, config.Error(400, "参数错误"))
		return
	}

	if err := h.svc.Create(&item); err != nil {
		c.JSON(http.StatusInternalServerError, config.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, config.Success(item))
}

func (h *IndustryStandardHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, config.Error(400, "ID无效"))
		return
	}

	var item model.IndustryStandard
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

func (h *IndustryStandardHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, config.Error(400, "ID无效"))
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, config.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, config.Success(nil))
}

func (h *IndustryStandardHandler) Check(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, config.Error(400, "ID无效"))
		return
	}

	var req struct {
		CheckStatus int `json:"check_status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, config.Error(400, "参数错误"))
		return
	}

	if err := h.svc.UpdateCheckStatus(uint(id), req.CheckStatus); err != nil {
		c.JSON(http.StatusInternalServerError, config.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, config.Success(nil))
}

func (h *IndustryStandardHandler) Import(c *gin.Context) {
	var items []model.IndustryStandard
	if err := c.ShouldBindJSON(&items); err != nil {
		c.JSON(http.StatusBadRequest, config.Error(400, "参数错误"))
		return
	}

	if err := h.svc.BatchImport(items); err != nil {
		c.JSON(http.StatusInternalServerError, config.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, config.Success(gin.H{"count": len(items)}))
}