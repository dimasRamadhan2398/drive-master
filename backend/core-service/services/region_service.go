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
	cacheSvc   ICacheService
}

func NewRegionService(regionRepo repositories.IRegionRepository, cacheSvc ICacheService) *RegionService {
	return &RegionService{
		regionRepo: regionRepo,
		cacheSvc:   cacheSvc,
	}
}

func (s *RegionService) GetAllProvinces(ctx context.Context) ([]models.Province, error) {
	cacheKey := "provinces:all"

	// Try to get from cache first
	var provinces []models.Province
	if err := s.cacheSvc.Get(ctx, cacheKey, &provinces); err == nil {
		return provinces, nil
	}

	// Get from repository
	provinces, err := s.regionRepo.GetProvinces(ctx)
	if err != nil {
		return nil, err
	}

	// Cache the result for 1 hour
	_ = s.cacheSvc.Set(ctx, cacheKey, provinces, 3600)

	return provinces, nil
}

func (s *RegionService) GetRegenciesByProvince(ctx context.Context, province string) ([]models.Regency, error) {
	cacheKey := "regencies:" + province

	// Try to get from cache first
	var regencies []models.Regency
	if err := s.cacheSvc.Get(ctx, cacheKey, &regencies); err == nil {
		return regencies, nil
	}

	// Get from repository
	regencies, err := s.regionRepo.GetRegencies(ctx, province)
	if err != nil {
		return nil, err
	}

	// Cache the result for 1 hour
	_ = s.cacheSvc.Set(ctx, cacheKey, regencies, 3600)

	return regencies, nil
}

func (s *RegionService) GetDistrictsByRegency(ctx context.Context, province, regency string) ([]models.District, error) {
	cacheKey := "districts:" + province + ":" + regency

	// Try to get from cache first
	var districts []models.District
	if err := s.cacheSvc.Get(ctx, cacheKey, &districts); err == nil {
		return districts, nil
	}

	// Get from repository
	districts, err := s.regionRepo.GetDistricts(ctx, province, regency)
	if err != nil {
		return nil, err
	}

	// Cache the result for 1 hour
	_ = s.cacheSvc.Set(ctx, cacheKey, districts, 3600)

	return districts, nil
}