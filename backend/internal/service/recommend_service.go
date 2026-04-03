package service

import (
	"policy-radar/internal/model"
	"policy-radar/internal/repository"
)

type RecommendService struct {
	repo *repository.RecommendRepo
}

func NewRecommendService(repo *repository.RecommendRepo) *RecommendService {
	return &RecommendService{repo: repo}
}

func (s *RecommendService) List(query *repository.RecommendQuery) ([]model.PolicyRecommend, int64, error) {
	return s.repo.List(query)
}

func (s *RecommendService) Create(item *model.PolicyRecommend) error {
	return s.repo.Create(item)
}

func (s *RecommendService) GetByID(id uint) (*model.PolicyRecommend, error) {
	return s.repo.GetByID(id)
}

func (s *RecommendService) GetByPolicy(policyID uint, policyType string) ([]model.PolicyRecommend, error) {
	return s.repo.GetByPolicy(policyID, policyType)
}

func (s *RecommendService) ListAll() ([]model.PolicyRecommend, error) {
	return s.repo.ListAll()
}

func (s *RecommendService) Update(item *model.PolicyRecommend) error {
	return s.repo.Update(item)
}

func (s *RecommendService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *RecommendService) Upsert(item *model.PolicyRecommend) error {
	return s.repo.Upsert(item)
}