package services

import (
	"context"
	"core-service/models"
	"core-service/pkg/kafka"
	"core-service/repositories"

	"github.com/google/uuid"
)

type ICarService interface {
	CreateCar(ctx context.Context, car *models.Car) error
	GetCarByID(ctx context.Context, id uuid.UUID) (*models.Car, error)
	GetAllCars(ctx context.Context) ([]models.Car, error)
	GetCarsByStatus(ctx context.Context, status models.CarStatus) ([]models.Car, error)
	GetCarsByTransmission(ctx context.Context, transmission models.TransmissionType) ([]models.Car, error)
	UpdateCar(ctx context.Context, car *models.Car) error
	DeleteCar(ctx context.Context, car *models.Car) error
	CountCars(ctx context.Context) (int64, error)
}

type CarService struct {
	carRepo        repositories.ICarRepository
	eventPublisher *kafka.EventPublisher
}

func NewCarService(carRepo repositories.ICarRepository, eventPublisher *kafka.EventPublisher) ICarService {
	return &CarService{
		carRepo:        carRepo,
		eventPublisher: eventPublisher,
	}
}

// CreateCar creates a new car
func (s *CarService) CreateCar(ctx context.Context, car *models.Car) error {
	if err := s.carRepo.Create(ctx, car); err != nil {
		return err
	}

	// Publish event (async to not block response)
	if s.eventPublisher != nil {
		go s.eventPublisher.PublishCarCreated(context.Background(), car)
	}

	return nil
}

// GetCarByID retrieves a car by ID
func (s *CarService) GetCarByID(ctx context.Context, id uuid.UUID) (*models.Car, error) {
	return s.carRepo.FindByID(ctx, id)
}

// GetAllCars retrieves all cars
func (s *CarService) GetAllCars(ctx context.Context) ([]models.Car, error) {
	return s.carRepo.FindAll(ctx)
}

// GetCarsByStatus retrieves cars by status
func (s *CarService) GetCarsByStatus(ctx context.Context, status models.CarStatus) ([]models.Car, error) {
	return s.carRepo.FindByStatus(ctx, status)
}

// GetCarsByTransmission retrieves cars by transmission type
func (s *CarService) GetCarsByTransmission(ctx context.Context, transmission models.TransmissionType) ([]models.Car, error) {
	return s.carRepo.FindByTransmission(ctx, transmission)
}

// UpdateCar updates a car
func (s *CarService) UpdateCar(ctx context.Context, car *models.Car) error {
	if err := s.carRepo.Update(ctx, car); err != nil {
		return err
	}

	// Publish event (async to not block response)
	if s.eventPublisher != nil {
		go s.eventPublisher.PublishCarUpdated(context.Background(), car)
	}

	return nil
}

// DeleteCar deletes a car
func (s *CarService) DeleteCar(ctx context.Context, car *models.Car) error {
	carID := car.ID.String()

	if err := s.carRepo.Delete(ctx, car); err != nil {
		return err
	}

	// Publish event (async to not block response)
	if s.eventPublisher != nil {
		go s.eventPublisher.PublishCarDeleted(context.Background(), carID)
	}

	return nil
}

// CountCars returns the total number of cars
func (s *CarService) CountCars(ctx context.Context) (int64, error) {
	return s.carRepo.Count(ctx)
}