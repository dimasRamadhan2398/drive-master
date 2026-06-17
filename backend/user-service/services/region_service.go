package services

import (
	"context"
	"fmt"

	"user-service/clients/region"
	"user-service/models"
)

// RegionService provides region lookups from core-service
type RegionService struct {
	regionClient *region.Client
}

func NewRegionService(regionClient *region.Client) *RegionService {
	return &RegionService{
		regionClient: regionClient,
	}
}

// GetAllProvinces returns all provinces from core-service
func (s *RegionService) GetAllProvinces(ctx context.Context) ([]region.Province, error) {
	return s.regionClient.GetAllProvinces(ctx)
}

// GetRegenciesByProvince returns regencies for a given province ID
func (s *RegionService) GetRegenciesByProvince(ctx context.Context, provinceID string) ([]region.Regency, error) {
	return s.regionClient.GetRegenciesByProvince(ctx, provinceID)
}

// GetDistrictsByRegency returns districts for a given regency and province ID
func (s *RegionService) GetDistrictsByRegency(ctx context.Context, provinceID, regencyID string) ([]region.District, error) {
	return s.regionClient.GetDistrictsByRegency(ctx, provinceID, regencyID)
}

// FindAreaByID looks up an area by type and ID
func (s *RegionService) FindAreaByID(ctx context.Context, areaType, areaID string) (string, error) {
	switch areaType {
	case "province":
		provinces, err := s.GetAllProvinces(ctx)
		if err != nil {
			return "", err
		}
		for _, p := range provinces {
			if p.ID == areaID {
				return p.Name, nil
			}
		}
		return "", fmt.Errorf("province not found: %s", areaID)

	case "regency":
		// Need to search through all provinces to find the regency
		provinces, err := s.GetAllProvinces(ctx)
		if err != nil {
			return "", err
		}
		for _, p := range provinces {
			regencies, err := s.GetRegenciesByProvince(ctx, p.ID)
			if err != nil {
				continue
			}
			for _, r := range regencies {
				if r.ID == areaID {
					return r.Name, nil
				}
			}
		}
		return "", fmt.Errorf("regency not found: %s", areaID)

	case "district":
		// Need to search all provinces and regencies to find the district
		provinces, err := s.GetAllProvinces(ctx)
		if err != nil {
			return "", err
		}
		for _, p := range provinces {
			regencies, err := s.GetRegenciesByProvince(ctx, p.ID)
			if err != nil {
				continue
			}
			for _, r := range regencies {
				districts, err := s.GetDistrictsByRegency(ctx, p.ID, r.ID)
				if err != nil {
					continue
				}
				for _, d := range districts {
					if d.ID == areaID {
						return d.Name, nil
					}
				}
			}
		}
		return "", fmt.Errorf("district not found: %s", areaID)

	default:
		return "", fmt.Errorf("invalid area type: %s", areaType)
	}
}

// FindAreaByIDWithTypes looks up an area by type and ID (using models types)
func (s *RegionService) FindAreaByIDWithTypes(ctx context.Context, areaType models.AreaType, areaID uint) (string, error) {
	return s.FindAreaByID(ctx, string(areaType), fmt.Sprintf("%d", areaID))
}

// IRegionService interface for dependency injection
type IRegionService interface {
	GetAllProvinces(ctx context.Context) ([]region.Province, error)
	GetRegenciesByProvince(ctx context.Context, provinceID string) ([]region.Regency, error)
	GetDistrictsByRegency(ctx context.Context, provinceID, regencyID string) ([]region.District, error)
	FindAreaByID(ctx context.Context, areaType, areaID string) (string, error)
	FindAreaByIDWithTypes(ctx context.Context, areaType models.AreaType, areaID uint) (string, error)
}
