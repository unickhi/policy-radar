package service

import (
	"policy-radar/internal/config"
	"policy-radar/internal/model"
	"policy-radar/internal/repository"
)

type NationalStandardService struct {
	repo *repository.NationalStandardRepo
}

func NewNationalStandardService(repo *repository.NationalStandardRepo) *NationalStandardService {
	return &NationalStandardService{repo: repo}
}

func (s *NationalStandardService) List(p *config.Pagination) ([]model.NationalStandard, int64, error) {
	return s.repo.List(p)
}

func (s *NationalStandardService) GetByID(id uint) (*model.NationalStandard, error) {
	return s.repo.GetByID(id)
}

func (s *NationalStandardService) Create(item *model.NationalStandard) error {
	item.StandardType = "国标"
	return s.repo.Create(item)
}

func (s *NationalStandardService) Update(item *model.NationalStandard) error {
	return s.repo.Update(item)
}

func (s *NationalStandardService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *NationalStandardService) UpdateCheckStatus(id uint, status int) error {
	return s.repo.UpdateCheckStatus(id, status)
}

func (s *NationalStandardService) BatchImport(items []model.NationalStandard) error {
	for i := range items {
		items[i].StandardType = "国标"
	}
	return s.repo.BatchCreate(items)
}