package repository

import (
	"policy-radar/internal/config"
	"policy-radar/internal/model"

	"gorm.io/gorm"
)

type LocalStandardRepo struct {
	db *gorm.DB
}

func NewLocalStandardRepo(db *gorm.DB) *LocalStandardRepo {
	return &LocalStandardRepo{db: db}
}

func (r *LocalStandardRepo) List(p *config.Pagination) ([]model.LocalStandard, int64, error) {
	var list []model.LocalStandard
	var total int64

	query := r.db.Model(&model.LocalStandard{})

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

func (r *LocalStandardRepo) GetByID(id uint) (*model.LocalStandard, error) {
	var item model.LocalStandard
	err := r.db.First(&item, id).Error
	return &item, err
}

func (r *LocalStandardRepo) Create(item *model.LocalStandard) error {
	return r.db.Create(item).Error
}

func (r *LocalStandardRepo) Update(item *model.LocalStandard) error {
	return r.db.Save(item).Error
}

func (r *LocalStandardRepo) Delete(id uint) error {
	return r.db.Delete(&model.LocalStandard{}, id).Error
}

func (r *LocalStandardRepo) UpdateCheckStatus(id uint, status int) error {
	return r.db.Model(&model.LocalStandard{}).Where("id = ?", id).Update("check_status", status).Error
}

func (r *LocalStandardRepo) BatchCreate(items []model.LocalStandard) error {
	return r.db.CreateInBatches(items, 100).Error
}

func (r *LocalStandardRepo) CountByCategory(categoryID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.LocalStandard{}).Where("category_id = ?", categoryID).Count(&count).Error
	return count, err
}

func (r *LocalStandardRepo) GetByStandardNo(standardNo string) (*model.LocalStandard, error) {
	var item model.LocalStandard
	err := r.db.Where("standard_no = ?", standardNo).First(&item).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &item, err
}

func (r *LocalStandardRepo) Upsert(item *model.LocalStandard) error {
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