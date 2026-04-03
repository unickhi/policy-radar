package service

import (
	"policy-radar/internal/config"
	"policy-radar/internal/model"
	"policy-radar/internal/repository"
)

type LocalStandardService struct {
	repo *repository.LocalStandardRepo
}

func NewLocalStandardService(repo *repository.LocalStandardRepo) *LocalStandardService {
	return &LocalStandardService{repo: repo}
}

func (s *LocalStandardService) List(p *config.Pagination) ([]model.LocalStandard, int64, error) {
	return s.repo.List(p)
}

func (s *LocalStandardService) GetByID(id uint) (*model.LocalStandard, error) {
	return s.repo.GetByID(id)
}

func (s *LocalStandardService) Create(item *model.LocalStandard) error {
	item.StandardType = "地标"
	return s.repo.Create(item)
}

func (s *LocalStandardService) Update(item *model.LocalStandard) error {
	return s.repo.Update(item)
}

func (s *LocalStandardService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *LocalStandardService) UpdateCheckStatus(id uint, status int) error {
	return s.repo.UpdateCheckStatus(id, status)
}

func (s *LocalStandardService) BatchImport(items []model.LocalStandard) error {
	for i := range items {
		items[i].StandardType = "地标"
	}
	return s.repo.BatchCreate(items)
}