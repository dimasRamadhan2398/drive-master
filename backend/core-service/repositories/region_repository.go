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
	options := &base.QueryOptions{
		Limit: 40,
	}
	if err := r.BaseRepository.FindMany(&models.District{}, &districts, options); err != nil {
		return nil, err
	}
	return districts, nil	
}

// GetProvinces implements [IRegionRepository].
func (r *RegionRepository) GetProvinces() ([]models.Province, error) {
	var provinces []models.Province
	options := &base.QueryOptions{
		Limit: 40,
	}
	if err := r.BaseRepository.FindMany(&models.Province{}, &provinces, options); err != nil {
		return nil, err
	}
	return provinces, nil	
}

// GetRegencies implements [IRegionRepository].
func (r *RegionRepository) GetRegencies(province string) ([]models.Regency, error) {
	var regencies []models.Regency
	options := &base.QueryOptions{
		Limit: 40,
	}
	if err := r.BaseRepository.FindMany(&models.Regency{}, &regencies, options); err != nil {
		return nil, err
	}
	return regencies, nil
}

func NewRegionRepository(db *base.BaseRepository) IRegionRepository {
	return &RegionRepository{BaseRepository: db}
}
