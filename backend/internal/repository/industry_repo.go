package repository

import (
	"policy-radar/internal/config"
	"policy-radar/internal/model"

	"gorm.io/gorm"
)

type IndustryStandardRepo struct {
	db *gorm.DB
}

func NewIndustryStandardRepo(db *gorm.DB) *IndustryStandardRepo {
	return &IndustryStandardRepo{db: db}
}

func (r *IndustryStandardRepo) List(p *config.Pagination) ([]model.IndustryStandard, int64, error) {
	var list []model.IndustryStandard
	var total int64

	query := r.db.Model(&model.IndustryStandard{})

	if p.Keyword != "" {
		query = query.Where("standard_no LIKE ? OR standard_name LIKE ?", "%"+p.Keyword+"%", "%"+p.Keyword+"%")
	}
	if p.Status != "" {
		query = query.Where("status = ?", p.Status)
	}
	if p.CheckStatus > 0 {
		query = query.Where("check_status = ?", p.CheckStatus)
	}
	if p.CategoryID > 0 {
		query = query.Where("category_id = ?", p.CategoryID)
	}

	query.Count(&total)
	err := query.Order("id DESC").Offset(p.GetOffset()).Limit(p.GetLimit()).Find(&list).Error
	return list, total, err
}

func (r *IndustryStandardRepo) GetByID(id uint) (*model.IndustryStandard, error) {
	var item model.IndustryStandard
	err := r.db.First(&item, id).Error
	return &item, err
}

func (r *IndustryStandardRepo) Create(item *model.IndustryStandard) error {
	return r.db.Create(item).Error
}

func (r *IndustryStandardRepo) Update(item *model.IndustryStandard) error {
	return r.db.Save(item).Error
}

func (r *IndustryStandardRepo) Delete(id uint) error {
	return r.db.Delete(&model.IndustryStandard{}, id).Error
}

func (r *IndustryStandardRepo) UpdateCheckStatus(id uint, status int) error {
	return r.db.Model(&model.IndustryStandard{}).Where("id = ?", id).Update("check_status", status).Error
}

func (r *IndustryStandardRepo) BatchCreate(items []model.IndustryStandard) error {
	return r.db.CreateInBatches(items, 100).Error
}

func (r *IndustryStandardRepo) CountByCategory(categoryID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.IndustryStandard{}).Where("category_id = ?", categoryID).Count(&count).Error
	return count, err
}

func (r *IndustryStandardRepo) GetByStandardNo(standardNo string) (*model.IndustryStandard, error) {
	var item model.IndustryStandard
	err := r.db.Where("standard_no = ?", standardNo).First(&item).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &item, err
}

func (r *IndustryStandardRepo) Upsert(item *model.IndustryStandard) error {
	existing, err := r.GetByStandardNo(item.StandardNo)
	if err != nil {
		return err
	}
	if existing != nil {
		item.ID = existing.ID
		item.CreatedAt = existing.CreatedAt
		return r.db.Save(item).Error
	}
	return r.Create(item)
}