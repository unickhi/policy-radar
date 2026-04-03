package service

import (
	"policy-radar/internal/config"
	"policy-radar/internal/model"
	"policy-radar/internal/repository"
)

type IndustryStandardService struct {
	repo *repository.IndustryStandardRepo
}

func NewIndustryStandardService(repo *repository.IndustryStandardRepo) *IndustryStandardService {
	return &IndustryStandardService{repo: repo}
}

func (s *IndustryStandardService) List(p *config.Pagination) ([]model.IndustryStandard, int64, error) {
	return s.repo.List(p)
}

func (s *IndustryStandardService) GetByID(id uint) (*model.IndustryStandard, error) {
	return s.repo.GetByID(id)
}

func (s *IndustryStandardService) Create(item *model.IndustryStandard) error {
	item.StandardType = "行标"
	return s.repo.Create(item)
}

func (s *IndustryStandardService) Update(item *model.IndustryStandard) error {
	return s.repo.Update(item)
}

func (s *IndustryStandardService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *IndustryStandardService) UpdateCheckStatus(id uint, status int) error {
	return s.repo.UpdateCheckStatus(id, status)
}

func (s *IndustryStandardService) BatchImport(items []model.IndustryStandard) error {
	for i := range items {
		items[i].StandardType = "行标"
	}
	return s.repo.BatchCreate(items)
}