package repositories

import (
	"context"
	"core-service/models"
	"core-service/pkg/base"

	"github.com/google/uuid"
)

type ICarRepository interface {
	Create(ctx context.Context, car *models.Car) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Car, error)
	FindAll(ctx context.Context) ([]models.Car, error)
	FindByStatus(ctx context.Context, status models.CarStatus) ([]models.Car, error)
	FindByTransmission(ctx context.Context, transmission models.TransmissionType) ([]models.Car, error)
	Update(ctx context.Context, car *models.Car) error
	Delete(ctx context.Context, car *models.Car) error
	Count(ctx context.Context) (int64, error)
}

type CarRepository struct {
	*base.BaseRepository
}

func NewCarRepository(baseRepo *base.BaseRepository) ICarRepository {
	return &CarRepository{BaseRepository: baseRepo}
}

// Create creates a new car
func (r *CarRepository) Create(ctx context.Context, car *models.Car) error {
	return r.BaseRepository.Create(car)
}

// FindByID finds a car by ID
func (r *CarRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Car, error) {
	var car models.Car
	if err := r.BaseRepository.FindByID(&car, id); err != nil {
		return nil, err
	}
	return &car, nil
}

// FindAll retrieves all cars
func (r *CarRepository) FindAll(ctx context.Context) ([]models.Car, error) {
	var cars []models.Car
	opts := base.NewQueryOptions()
	if err := r.BaseRepository.FindMany(&models.Car{}, &cars, opts); err != nil {
		return nil, err
	}
	return cars, nil
}

// FindByStatus retrieves cars by status
func (r *CarRepository) FindByStatus(ctx context.Context, status models.CarStatus) ([]models.Car, error) {
	var cars []models.Car
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"status": status})
	if err := r.BaseRepository.FindMany(&models.Car{}, &cars, opts); err != nil {
		return nil, err
	}
	return cars, nil
}

// FindByTransmission retrieves cars by transmission type
func (r *CarRepository) FindByTransmission(ctx context.Context, transmission models.TransmissionType) ([]models.Car, error) {
	var cars []models.Car
	opts := base.NewQueryOptions().
		WithWhere(map[string]any{"transmission": transmission})
	if err := r.BaseRepository.FindMany(&models.Car{}, &cars, opts); err != nil {
		return nil, err
	}
	return cars, nil
}

// Update updates a car
func (r *CarRepository) Update(ctx context.Context, car *models.Car) error {
	return r.BaseRepository.Update(car)
}

// Delete deletes a car
func (r *CarRepository) Delete(ctx context.Context, car *models.Car) error {
	return r.BaseRepository.Delete(car)
}

// Count returns the total number of cars
func (r *CarRepository) Count(ctx context.Context) (int64, error) {
	return r.BaseRepository.Count(&models.Car{}, nil)
}