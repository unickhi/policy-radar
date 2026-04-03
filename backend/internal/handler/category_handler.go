package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"policy-radar/internal/config"
	"policy-radar/internal/model"
	"policy-radar/internal/service"
)

type CategoryHandler struct {
	svc *service.CategoryService
}

func NewCategoryHandler(svc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

func (h *CategoryHandler) List(c *gin.Context) {
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

func (h *CategoryHandler) Get(c *gin.Context) {
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

func (h *CategoryHandler) Create(c *gin.Context) {
	var item model.PolicyCategory
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

func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, config.Error(400, "ID无效"))
		return
	}

	var item model.PolicyCategory
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

func (h *CategoryHandler) Delete(c *gin.Context) {
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

func (h *CategoryHandler) Count(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, config.Error(400, "ID无效"))
		return
	}

	count, err := h.svc.GetCountByCategory(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, config.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, config.Success(count))
}