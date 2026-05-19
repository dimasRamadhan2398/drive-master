package services

import (
	"context"
	"core-service/models"
	"core-service/repositories"
)

type IRegionService interface {
	GetAllProvinces(ctx context.Context) ([]models.Province, error)
	GetRegenciesByProvince(ctx context.Context, province string) ([]models.Regency, error)
	GetDistrictsByRegency(ctx context.Context, province, regency string) ([]models.District, error)
}

type RegionService struct {
	regionRepo repositories.IRegionRepository
}

func NewRegionService(regionRepo repositories.IRegionRepository) *RegionService {
	return &RegionService{
		regionRepo: regionRepo,
	}
}

func (s *RegionService) GetAllProvinces(ctx context.Context) ([]models.Province, error) {
	return s.regionRepo.GetProvinces(ctx)
}

func (s *RegionService) GetRegenciesByProvince(ctx context.Context, province string) ([]models.Regency, error) {
	return s.regionRepo.GetRegencies(ctx, province)
}

func (s *RegionService) GetDistrictsByRegency(ctx context.Context, province, regency string) ([]models.District, error) {
	return s.regionRepo.GetDistricts(ctx, province, regency)
}