package service

import (
	"policy-radar/internal/config"
	"policy-radar/internal/model"
	"policy-radar/internal/repository"
)

type CategoryService struct {
	repo        *repository.CategoryRepo
	nationalRepo *repository.NationalStandardRepo
	industryRepo *repository.IndustryStandardRepo
	localRepo   *repository.LocalStandardRepo
}

func NewCategoryService(repo *repository.CategoryRepo, nationalRepo *repository.NationalStandardRepo, industryRepo *repository.IndustryStandardRepo, localRepo *repository.LocalStandardRepo) *CategoryService {
	return &CategoryService{
		repo:        repo,
		nationalRepo: nationalRepo,
		industryRepo: industryRepo,
		localRepo:   localRepo,
	}
}

func (s *CategoryService) List(p *config.Pagination) ([]model.PolicyCategory, int64, error) {
	return s.repo.List(p)
}

func (s *CategoryService) GetByID(id uint) (*model.PolicyCategory, error) {
	return s.repo.GetByID(id)
}

func (s *CategoryService) Create(item *model.PolicyCategory) error {
	return s.repo.Create(item)
}

func (s *CategoryService) Update(item *model.PolicyCategory) error {
	return s.repo.Update(item)
}

func (s *CategoryService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *CategoryService) GetCountByCategory(categoryID uint) (map[string]int64, error) {
	national, err := s.nationalRepo.CountByCategory(categoryID)
	if err != nil {
		return nil, err
	}
	industry, err := s.industryRepo.CountByCategory(categoryID)
	if err != nil {
		return nil, err
	}
	local, err := s.localRepo.CountByCategory(categoryID)
	if err != nil {
		return nil, err
	}
	return map[string]int64{
		"national": national,
		"industry": industry,
		"local":    local,
		"total":    national + industry + local,
	}, nil
}

func (s *CategoryService) GetAll() ([]model.PolicyCategory, error) {
	return s.repo.GetAll()
}