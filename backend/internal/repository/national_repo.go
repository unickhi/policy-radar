package repository

import (
	"policy-radar/internal/config"
	"policy-radar/internal/model"

	"gorm.io/gorm"
)

type NationalStandardRepo struct {
	db *gorm.DB
}

func NewNationalStandardRepo(db *gorm.DB) *NationalStandardRepo {
	return &NationalStandardRepo{db: db}
}

func (r *NationalStandardRepo) List(p *config.Pagination) ([]model.NationalStandard, int64, error) {
	var list []model.NationalStandard
	var total int64

	query := r.db.Model(&model.NationalStandard{})

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

func (r *NationalStandardRepo) GetByID(id uint) (*model.NationalStandard, error) {
	var item model.NationalStandard
	err := r.db.First(&item, id).Error
	return &item, err
}

func (r *NationalStandardRepo) Create(item *model.NationalStandard) error {
	return r.db.Create(item).Error
}

func (r *NationalStandardRepo) Update(item *model.NationalStandard) error {
	return r.db.Save(item).Error
}

func (r *NationalStandardRepo) Delete(id uint) error {
	return r.db.Delete(&model.NationalStandard{}, id).Error
}

func (r *NationalStandardRepo) UpdateCheckStatus(id uint, status int) error {
	return r.db.Model(&model.NationalStandard{}).Where("id = ?", id).Update("check_status", status).Error
}

func (r *NationalStandardRepo) BatchCreate(items []model.NationalStandard) error {
	return r.db.CreateInBatches(items, 100).Error
}

func (r *NationalStandardRepo) CountByCategory(categoryID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.NationalStandard{}).Where("category_id = ?", categoryID).Count(&count).Error
	return count, err
}

func (r *NationalStandardRepo) GetByStandardNo(standardNo string) (*model.NationalStandard, error) {
	var item model.NationalStandard
	err := r.db.Where("standard_no = ?", standardNo).First(&item).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &item, err
}

func (r *NationalStandardRepo) GetByStandardNoAndName(standardNo, standardName string) (*model.NationalStandard, error) {
	var item model.NationalStandard
	err := r.db.Where("standard_no = ? AND standard_name = ?", standardNo, standardName).First(&item).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &item, err
}

func (r *NationalStandardRepo) Upsert(item *model.NationalStandard) error {
	existing, err := r.GetByStandardNoAndName(item.StandardNo, item.StandardName)
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