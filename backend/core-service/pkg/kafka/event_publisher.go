package kafka

import (
	"context"
	"time"

	"core-service/models"
)

// EventType constants
const (
	EventCarCreated    = "car.created"
	EventCarUpdated    = "car.updated"
	EventCarDeleted    = "car.deleted"
	EventPackageCreated = "package.created"
	EventPackageUpdated = "package.updated"
	EventPackageDeleted = "package.deleted"
)

// EventPublisher handles publishing domain events
type EventPublisher struct {
	producer *KafkaProducer
}

// NewEventPublisher creates a new event publisher
func NewEventPublisher(producer *KafkaProducer) *EventPublisher {
	return &EventPublisher{producer: producer}
}

// PublishCarCreated publishes a car created event
func (e *EventPublisher) PublishCarCreated(ctx context.Context, car *models.Car) error {
	event := models.CarCreatedEvent{
		CarID:         car.ID.String(),
		Brand:         car.Brand,
		Model:         car.Model,
		Year:          car.Year,
		LicensePlate:  car.LicensePlate,
		Transmission:  string(car.Transmission),
		CreatedAt:     time.Now().Format(time.RFC3339),
	}
	return e.producer.Publish(ctx, EventCarCreated, event)
}

// PublishCarUpdated publishes a car updated event
func (e *EventPublisher) PublishCarUpdated(ctx context.Context, car *models.Car) error {
	event := models.CarUpdatedEvent{
		CarID:      car.ID.String(),
		Brand:      car.Brand,
		Model:      car.Model,
		Year:       car.Year,
		Status:     string(car.Status),
		UpdatedAt:  time.Now().Format(time.RFC3339),
	}
	return e.producer.Publish(ctx, EventCarUpdated, event)
}

// PublishCarDeleted publishes a car deleted event
func (e *EventPublisher) PublishCarDeleted(ctx context.Context, carID string) error {
	event := models.CarDeletedEvent{
		CarID: carID,
	}
	return e.producer.Publish(ctx, EventCarDeleted, event)
}

// PublishPackageCreated publishes a package created event
func (e *EventPublisher) PublishPackageCreated(ctx context.Context, pkg *models.Package) error {
	event := models.PackageCreatedEvent{
		PackageID:   pkg.ID.String(),
		Name:        pkg.Name,
		PackageType: string(pkg.PackageType),
		Price:       pkg.Price,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}
	return e.producer.Publish(ctx, EventPackageCreated, event)
}

// PublishPackageUpdated publishes a package updated event
func (e *EventPublisher) PublishPackageUpdated(ctx context.Context, pkg *models.Package) error {
	event := models.PackageUpdatedEvent{
		PackageID:   pkg.ID.String(),
		Name:        pkg.Name,
		PackageType: string(pkg.PackageType),
		Price:       pkg.Price,
		Status:      string(pkg.Status),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}
	return e.producer.Publish(ctx, EventPackageUpdated, event)
}

// PublishPackageDeleted publishes a package deleted event
func (e *EventPublisher) PublishPackageDeleted(ctx context.Context, packageID string) error {
	event := models.PackageDeletedEvent{
		PackageID: packageID,
	}
	return e.producer.Publish(ctx, EventPackageDeleted, event)
}