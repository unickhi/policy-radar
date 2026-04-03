package repository

import (
	"policy-radar/internal/model"

	"gorm.io/gorm"
)

type RecommendRepo struct {
	db *gorm.DB
}

func NewRecommendRepo(db *gorm.DB) *RecommendRepo {
	return &RecommendRepo{db: db}
}

type RecommendQuery struct {
	Keyword    string
	PolicyType string
	Page       int
	PageSize   int
}

func (r *RecommendRepo) List(query *RecommendQuery) ([]model.PolicyRecommend, int64, error) {
	var list []model.PolicyRecommend
	var total int64

	db := r.db.Model(&model.PolicyRecommend{})

	// 搜索条件
	if query.Keyword != "" {
		db = db.Where("title LIKE ? OR policy_name LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}
	if query.PolicyType != "" {
		db = db.Where("policy_type = ?", query.PolicyType)
	}

	// 统计总数
	db.Count(&total)

	// 分页查询
	offset := (query.Page - 1) * query.PageSize
	err := db.Order("sort ASC, id DESC").Offset(offset).Limit(query.PageSize).Find(&list).Error
	return list, total, err
}

func (r *RecommendRepo) Create(item *model.PolicyRecommend) error {
	return r.db.Create(item).Error
}

func (r *RecommendRepo) GetByID(id uint) (*model.PolicyRecommend, error) {
	var item model.PolicyRecommend
	err := r.db.First(&item, id).Error
	return &item, err
}

func (r *RecommendRepo) GetByPolicy(policyID uint, policyType string) ([]model.PolicyRecommend, error) {
	var list []model.PolicyRecommend
	err := r.db.Where("policy_id = ? AND policy_type = ?", policyID, policyType).Order("sort ASC").Find(&list).Error
	return list, err
}

func (r *RecommendRepo) ListAll() ([]model.PolicyRecommend, error) {
	var list []model.PolicyRecommend
	err := r.db.Order("sort ASC, id DESC").Find(&list).Error
	return list, err
}

func (r *RecommendRepo) Update(item *model.PolicyRecommend) error {
	return r.db.Save(item).Error
}

func (r *RecommendRepo) Delete(id uint) error {
	return r.db.Delete(&model.PolicyRecommend{}, id).Error
}

func (r *RecommendRepo) Upsert(item *model.PolicyRecommend) error {
	if item.ID > 0 {
		return r.db.Save(item).Error
	}
	return r.db.Create(item).Error
}