package repositories

import (
	"core-service/models"
	"core-service/pkg/base"
)

type IRegionRepository interface {
	GetProvinces() ([]models.Province, error)
	GetRegencies(province string) ([]models.Regency, error)
	GetDistricts(province, regency string) ([]models.District, error)
}

type RegionRepository struct {
	*base.BaseRepository
}

// GetDistricts implements [IRegionRepository].
func (r *RegionRepository) GetDistricts(province string, regency string) ([]models.District, error) {
	var districts []models.District
	if err := r.BaseRepository.FindMany(&models.District{}, &districts, nil); err != nil {
		return nil, err
	}
	return districts, nil	
}

// GetProvinces implements [IRegionRepository].
func (r *RegionRepository) GetProvinces() ([]models.Province, error) {
	var provinces []models.Province
	if err := r.BaseRepository.FindMany(&models.Province{}, &provinces, nil); err != nil {
		return nil, err
	}
	return provinces, nil	
}

// GetRegencies implements [IRegionRepository].
func (r *RegionRepository) GetRegencies(province string) ([]models.Regency, error) {
	var regencies []models.Regency
	if err := r.BaseRepository.FindMany(&models.Regency{}, &regencies, nil); err != nil {
		return nil, err
	}
	return regencies, nil
}

func NewRegionRepository(db *base.BaseRepository) IRegionRepository {
	return &RegionRepository{BaseRepository: db}
}
