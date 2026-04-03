package repository

import (
	"policy-radar/internal/model"

	"gorm.io/gorm"
)

type CrawlerLogRepo struct {
	db *gorm.DB
}

func NewCrawlerLogRepo(db *gorm.DB) *CrawlerLogRepo {
	return &CrawlerLogRepo{db: db}
}

func (r *CrawlerLogRepo) Create(log *model.CrawlerLog) error {
	return r.db.Create(log).Error
}

func (r *CrawlerLogRepo) List(limit int) ([]model.CrawlerLog, error) {
	var logs []model.CrawlerLog
	err := r.db.Order("id DESC").Limit(limit).Find(&logs).Error
	return logs, err
}