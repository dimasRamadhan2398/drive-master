package services

import (
	"core-service/models"
	"core-service/repositories"
)

type RegionService struct {
	regionRepo repositories.IRegionRepository
}

type IRegionService interface {
	GetAllProvinces() ([]models.Province, error)
	GetRegenciesByProvince(province string) ([]models.Regency, error)
	GetDistrictsByRegency(province, regency string) ([]models.District, error)
}

func NewRegionService(regionRepo repositories.IRegionRepository) *RegionService {
	return &RegionService{
		regionRepo: regionRepo,
	}
}

func (s *RegionService) GetAllProvinces() ([]models.Province, error) {
	return s.regionRepo.GetProvinces()
}

func (s *RegionService) GetRegenciesByProvince(province string) ([]models.Regency, error) {
	return s.regionRepo.GetRegencies(province)
}

func (s *RegionService) GetDistrictsByRegency(province, regency string) ([]models.District, error) {
	return s.regionRepo.GetDistricts(province, regency)
}