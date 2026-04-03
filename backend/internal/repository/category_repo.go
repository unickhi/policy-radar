package repository

import (
	"policy-radar/internal/config"
	"policy-radar/internal/model"

	"gorm.io/gorm"
)

type CategoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) List(p *config.Pagination) ([]model.PolicyCategory, int64, error) {
	var list []model.PolicyCategory
	var total int64

	query := r.db.Model(&model.PolicyCategory{})

	if p.Keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+p.Keyword+"%", "%"+p.Keyword+"%")
	}

	query.Count(&total)
	err := query.Order("id DESC").Offset(p.GetOffset()).Limit(p.GetLimit()).Find(&list).Error
	return list, total, err
}

func (r *CategoryRepo) GetByID(id uint) (*model.PolicyCategory, error) {
	var item model.PolicyCategory
	err := r.db.First(&item, id).Error
	return &item, err
}

func (r *CategoryRepo) Create(item *model.PolicyCategory) error {
	return r.db.Create(item).Error
}

func (r *CategoryRepo) Update(item *model.PolicyCategory) error {
	return r.db.Save(item).Error
}

func (r *CategoryRepo) Delete(id uint) error {
	return r.db.Delete(&model.PolicyCategory{}, id).Error
}

func (r *CategoryRepo) GetAll() ([]model.PolicyCategory, error) {
	var list []model.PolicyCategory
	err := r.db.Find(&list).Error
	return list, err
}