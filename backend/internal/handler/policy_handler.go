package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"policy-radar/internal/config"
	"policy-radar/internal/service"
)

// PolicyHandler 前端展示用的政策查询处理器
type PolicyHandler struct {
	nationalSvc  *service.NationalStandardService
	industrySvc  *service.IndustryStandardService
	localSvc     *service.LocalStandardService
	recommendSvc *service.RecommendService
	categorySvc  *service.CategoryService
}

func NewPolicyHandler(
	nationalSvc *service.NationalStandardService,
	industrySvc *service.IndustryStandardService,
	localSvc *service.LocalStandardService,
	recommendSvc *service.RecommendService,
	categorySvc *service.CategoryService,
) *PolicyHandler {
	return &PolicyHandler{
		nationalSvc:  nationalSvc,
		industrySvc:  industrySvc,
		localSvc:     localSvc,
		recommendSvc: recommendSvc,
		categorySvc:  categorySvc,
	}
}

func (h *PolicyHandler) ListAll(c *gin.Context) {
	var p config.Pagination
	if err := c.ShouldBindQuery(&p); err != nil {
		c.JSON(http.StatusBadRequest, config.Error(400, "参数错误"))
		return
	}

	// 获取三种类型的政策
	nationalList, nationalTotal, _ := h.nationalSvc.List(&p)
	industryList, industryTotal, _ := h.industrySvc.List(&p)
	localList, localTotal, _ := h.localSvc.List(&p)

	// 合并结果
	result := gin.H{
		"national": gin.H{
			"list":  nationalList,
			"total": nationalTotal,
		},
		"industry": gin.H{
			"list":  industryList,
			"total": industryTotal,
		},
		"local": gin.H{
			"list":  localList,
			"total": localTotal,
		},
	}

	c.JSON(http.StatusOK, config.Success(result))
}

// GetDetail 获取政策详情（前台H5使用）
func (h *PolicyHandler) GetDetail(c *gin.Context) {
	policyType := c.Param("type")
	idStr := c.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, config.Error(400, "ID无效"))
		return
	}

	var detail gin.H
	var categoryID uint
	var standardName string

	switch policyType {
	case "national":
		item, err := h.nationalSvc.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, config.Error(404, "记录不存在"))
			return
		}
		categoryID = item.CategoryID
		standardName = item.StandardName
		// 政策溯源使用link1
		sourceLink := item.Link1
		detail = gin.H{
			"id":             item.ID,
			"standard_no":    item.StandardNo,
			"standard_name":  item.StandardName,
			"english_name":   item.EnglishName,
			"publish_date":   item.PublishDate,
			"implement_date": item.ImplementDate,
			"status":         item.Status,
			"nature":         item.Nature,
			"category":       item.Category,
			"is_adopted":     item.IsAdopted,
			"ccs_code":       item.CCSCode,
			"ics_code":       item.ICSCode,
			"department":     item.Department,
			"technical_dept": item.TechnicalDept,
			"publisher":      item.Publisher,
			"description":    item.Description,
			"download_url":   item.DownloadURL,
			"category_id":    item.CategoryID,
			"source_link":    sourceLink,
			"policy_type":    "国标",
		}

	case "industry":
		item, err := h.industrySvc.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, config.Error(404, "记录不存在"))
			return
		}
		categoryID = item.CategoryID
		standardName = item.StandardName
		// 政策溯源使用detail_link
		sourceLink := item.DetailLink
		detail = gin.H{
			"id":              item.ID,
			"standard_no":     item.StandardNo,
			"standard_name":   item.StandardName,
			"publish_date":    item.PublishDate,
			"implement_date":  item.ImplementDate,
			"status":          item.Status,
			"revision_type":   item.RevisionType,
			"ccs_code":        item.CCSCode,
			"ics_code":        item.ICSCode,
			"technical_owner": item.TechnicalOwner,
			"approve_dept":    item.ApproveDept,
			"industry_class":  item.IndustryClass,
			"standard_class":  item.StandardClass,
			"replace_standard": item.ReplaceStandard,
			"download_url":    item.DownloadURL,
			"category_id":     item.CategoryID,
			"source_link":     sourceLink,
			"policy_type":     "行标",
		}

	case "local":
		item, err := h.localSvc.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, config.Error(404, "记录不存在"))
			return
		}
		categoryID = item.CategoryID
		standardName = item.StandardName
		// 政策溯源使用detail_link
		sourceLink := item.DetailLink
		detail = gin.H{
			"id":             item.ID,
			"standard_no":    item.StandardNo,
			"standard_name":  item.StandardName,
			"publish_date":   item.PublishDate,
			"implement_date": item.ImplementDate,
			"status":         item.Status,
			"nature":         item.Nature,
			"ccs_code":       item.CCSCode,
			"ics_code":       item.ICSCode,
			"department":     item.Department,
			"publisher":      item.Publisher,
			"description":    item.Description,
			"download_url":   item.DownloadURL,
			"category_id":    item.CategoryID,
			"source_link":    sourceLink,
			"policy_type":    "地标",
		}

	default:
		c.JSON(http.StatusBadRequest, config.Error(400, "无效的政策类型"))
		return
	}

	// 获取政策拆分内容
	policySplits, _ := h.recommendSvc.GetByPolicy(uint(id), policyType)
	detail["policy_splits"] = policySplits

	// 获取相关推荐（同分类优先，无分类则用关键词匹配）
	relatedList := h.getRelatedPolicies(policyType, categoryID, uint(id), standardName)
	detail["related"] = relatedList

	// 获取分类名称
	if categoryID > 0 {
		category, err := h.categorySvc.GetByID(categoryID)
		if err == nil {
			detail["category_name"] = category.Name
		}
	}

	c.JSON(http.StatusOK, config.Success(detail))
}

// getRelatedPolicies 获取相关推荐政策（同分类，跨类型，优先地标）
// 如果 categoryID = 0 表示没有分类，返回空
func (h *PolicyHandler) getRelatedPolicies(policyType string, categoryID uint, excludeID uint, standardName string) []gin.H {
	var result []gin.H

	// 没有分类，返回空
	if categoryID == 0 {
		return result
	}

	// 优先查找地标，然后是当前类型，最后是其他类型
	typeOrder := []string{"local", policyType}
	if policyType != "national" {
		typeOrder = append(typeOrder, "national")
	}
	if policyType != "industry" {
		typeOrder = append(typeOrder, "industry")
	}

	p := &config.Pagination{
		Page:       1,
		PageSize:   10,
		CategoryID: categoryID,
	}

	for _, t := range typeOrder {
		if len(result) >= 5 {
			break
		}

		var list []gin.H
		switch t {
		case "national":
			items, _, _ := h.nationalSvc.List(p)
			for _, item := range items {
				if t == policyType && item.ID == excludeID {
					continue
				}
				list = append(list, gin.H{
					"id":            item.ID,
					"type":          "national",
					"standard_no":   item.StandardNo,
					"standard_name": item.StandardName,
					"publish_date":  item.PublishDate,
					"status":        item.Status,
				})
			}
		case "industry":
			items, _, _ := h.industrySvc.List(p)
			for _, item := range items {
				if t == policyType && item.ID == excludeID {
					continue
				}
				list = append(list, gin.H{
					"id":            item.ID,
					"type":          "industry",
					"standard_no":   item.StandardNo,
					"standard_name": item.StandardName,
					"publish_date":  item.PublishDate,
					"status":        item.Status,
				})
			}
		case "local":
			items, _, _ := h.localSvc.List(p)
			for _, item := range items {
				if t == policyType && item.ID == excludeID {
					continue
				}
				list = append(list, gin.H{
					"id":            item.ID,
					"type":          "local",
					"standard_no":   item.StandardNo,
					"standard_name": item.StandardName,
					"publish_date":  item.PublishDate,
					"status":        item.Status,
				})
			}
		}

		for _, item := range list {
			if len(result) >= 5 {
				break
			}
			result = append(result, item)
		}
	}

	return result
}

// extractKeyword 从标准名称中提取关键词（已移除，不再需要）

func (h *PolicyHandler) Search(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, config.Error(400, "搜索关键词不能为空"))
		return
	}

	p := &config.Pagination{
		Keyword:  keyword,
		Page:     1,
		PageSize: 20,
	}

	// 在三种类型中搜索
	nationalList, _, _ := h.nationalSvc.List(p)
	industryList, _, _ := h.industrySvc.List(p)
	localList, _, _ := h.localSvc.List(p)

	// 收集所有分类ID
	categoryIDs := make(map[uint]bool)
	for _, item := range nationalList {
		if item.CategoryID > 0 {
			categoryIDs[item.CategoryID] = true
		}
	}
	for _, item := range industryList {
		if item.CategoryID > 0 {
			categoryIDs[item.CategoryID] = true
		}
	}
	for _, item := range localList {
		if item.CategoryID > 0 {
			categoryIDs[item.CategoryID] = true
		}
	}

	// 获取分类名称映射
	categoryNames := make(map[uint]string)
	for catID := range categoryIDs {
		cat, err := h.categorySvc.GetByID(catID)
		if err == nil {
			categoryNames[catID] = cat.Name
		}
	}

	// 转换为统一格式
	var results []gin.H
	for _, item := range nationalList {
		result := gin.H{
			"id":            item.ID,
			"type":          "national",
			"type_name":     "国标",
			"standard_no":   item.StandardNo,
			"standard_name": item.StandardName,
			"publish_date":  item.PublishDate,
			"status":        item.Status,
			"publisher":     item.Publisher,
			"category_id":   item.CategoryID,
		}
		if item.CategoryID > 0 && categoryNames[item.CategoryID] != "" {
			result["category_name"] = categoryNames[item.CategoryID]
		}
		results = append(results, result)
	}
	for _, item := range industryList {
		result := gin.H{
			"id":            item.ID,
			"type":          "industry",
			"type_name":     "行标",
			"standard_no":   item.StandardNo,
			"standard_name": item.StandardName,
			"publish_date":  item.PublishDate,
			"status":        item.Status,
			"publisher":     item.ApproveDept,
			"category_id":   item.CategoryID,
		}
		if item.CategoryID > 0 && categoryNames[item.CategoryID] != "" {
			result["category_name"] = categoryNames[item.CategoryID]
		}
		results = append(results, result)
	}
	for _, item := range localList {
		result := gin.H{
			"id":            item.ID,
			"type":          "local",
			"type_name":     "地标",
			"standard_no":   item.StandardNo,
			"standard_name": item.StandardName,
			"publish_date":  item.PublishDate,
			"status":        item.Status,
			"publisher":     item.Publisher,
			"category_id":   item.CategoryID,
		}
		if item.CategoryID > 0 && categoryNames[item.CategoryID] != "" {
			result["category_name"] = categoryNames[item.CategoryID]
		}
		results = append(results, result)
	}

	c.JSON(http.StatusOK, config.Success(results))
}

// ListByType 按类型获取列表（前台H5使用）
func (h *PolicyHandler) ListByType(c *gin.Context) {
	policyType := c.Param("type")

	var p config.Pagination
	if err := c.ShouldBindQuery(&p); err != nil {
		c.JSON(http.StatusBadRequest, config.Error(400, "参数错误"))
		return
	}

	var list interface{}
	var total int64

	switch policyType {
	case "national":
		list, total, _ = h.nationalSvc.List(&p)
	case "industry":
		list, total, _ = h.industrySvc.List(&p)
	case "local":
		list, total, _ = h.localSvc.List(&p)
	default:
		c.JSON(http.StatusBadRequest, config.Error(400, "无效的政策类型"))
		return
	}

	c.JSON(http.StatusOK, config.PageResponse{
		Code:    0,
		Message: "success",
		Data: &config.PageData{
			List:     list,
			Total:    total,
			Page:     p.Page,
			PageSize: p.PageSize,
		},
	})
}

// GetHomeData 获取首页数据（Banner + 各类型预览）
func (h *PolicyHandler) GetHomeData(c *gin.Context) {
	// 获取推荐政策作为Banner
	recommendList, _ := h.recommendSvc.ListAll()
	var banners []gin.H
	for i, r := range recommendList {
		if i >= 3 {
			break
		}
		banners = append(banners, gin.H{
			"id":          r.ID,
			"title":       r.Title,
			"content":     r.Content,
			"policy_id":   r.PolicyID,
			"policy_type": r.PolicyType,
		})
	}

	// 获取各类型预览（前3条）
	p := &config.Pagination{Page: 1, PageSize: 3}

	nationalList, nationalTotal, _ := h.nationalSvc.List(p)
	industryList, industryTotal, _ := h.industrySvc.List(p)
	localList, localTotal, _ := h.localSvc.List(p)

	c.JSON(http.StatusOK, config.Success(gin.H{
		"banners": banners,
		"national": gin.H{
			"list":  nationalList,
			"total": nationalTotal,
		},
		"industry": gin.H{
			"list":  industryList,
			"total": industryTotal,
		},
		"local": gin.H{
			"list":  localList,
			"total": localTotal,
		},
	}))
}