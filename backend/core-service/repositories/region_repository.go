package repositories

import (
	"context"
	"core-service/models"
	"core-service/pkg/base"
)

type IRegionRepository interface {
	GetProvinces(ctx context.Context) ([]models.Province, error)
	GetRegencies(ctx context.Context, province string) ([]models.Regency, error)
	GetDistricts(ctx context.Context, province, regency string) ([]models.District, error)
}

type RegionRepository struct {
	*base.BaseRepository
}

// GetDistricts implements [IRegionRepository].
func (r *RegionRepository) GetDistricts(ctx context.Context, province string, regency string) ([]models.District, error) {
	var districts []models.District
	options := &base.QueryOptions{
		Limit: 40,
	}
	if err := r.BaseRepository.FindMany(ctx, &models.District{}, &districts, options); err != nil {
		return nil, err
	}
	return districts, nil
}

// GetProvinces implements [IRegionRepository].
func (r *RegionRepository) GetProvinces(ctx context.Context) ([]models.Province, error) {
	var provinces []models.Province
	options := &base.QueryOptions{
		Limit: 40,
	}
	if err := r.BaseRepository.FindMany(ctx, &models.Province{}, &provinces, options); err != nil {
		return nil, err
	}
	return provinces, nil
}

// GetRegencies implements [IRegionRepository].
func (r *RegionRepository) GetRegencies(ctx context.Context, province string) ([]models.Regency, error) {
	var regencies []models.Regency
	options := &base.QueryOptions{
		Limit: 40,
	}
	if err := r.BaseRepository.FindMany(ctx, &models.Regency{}, &regencies, options); err != nil {
		return nil, err
	}
	return regencies, nil
}

func NewRegionRepository(db *base.BaseRepository) IRegionRepository {
	return &RegionRepository{BaseRepository: db}
}