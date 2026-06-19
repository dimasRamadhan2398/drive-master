package kafka

import (
	"context"
	"time"

	"core-service/models"
)

// EventType constants
const (
	// ========== CAR EVENTS ==========
	EventCarCreated    = "car.created"
	EventCarUpdated    = "car.updated"
	EventCarDeleted    = "car.deleted"

	// ========== PACKAGE EVENTS ==========
	EventPackageCreated = "package.created"
	EventPackageUpdated = "package.updated"
	EventPackageDeleted = "package.deleted"

	// ========== ARTICLE EVENTS ==========
	EventArticleCreated   = "article.created"
	EventArticleUpdated   = "article.updated"
	EventArticleDeleted   = "article.deleted"
	EventArticlePublished = "article.published"
	EventArticleArchived  = "article.archived"

	// ========== USER EVENTS (upstream - consumed from other services) ==========
	EventUserCreated = "user.created"
	EventUserUpdated = "user.updated"
	EventUserDeleted = "user.deleted"

	// ========== REGION EVENTS ==========
	EventRegionProvinceUpdated = "region.province.updated"
	EventRegionRegencyUpdated  = "region.regency.updated"
	EventRegionDistrictUpdated = "region.district.updated"

	// ========== SYSTEM EVENTS ==========
	EventCacheInvalidated = "cache.invalidated"
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

// ========== ARTICLE EVENTS ==========

// PublishArticleCreated publishes an article created event
func (e *EventPublisher) PublishArticleCreated(ctx context.Context, article *models.Article) error {
	event := models.ArticleCreatedEvent{
		ArticleID:   article.ID.String(),
		Title:       article.Title,
		Slug:        article.Slug,
		Status:      string(article.Status),
		AuthorID:    article.AuthorID.String(),
		IsFeatured:  article.IsFeatured,
		IsSpotlight: article.IsSpotlight,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}
	return e.producer.Publish(ctx, EventArticleCreated, event)
}

// PublishArticleUpdated publishes an article updated event
func (e *EventPublisher) PublishArticleUpdated(ctx context.Context, article *models.Article) error {
	event := models.ArticleUpdatedEvent{
		ArticleID:   article.ID.String(),
		Title:       article.Title,
		Slug:        article.Slug,
		Status:      string(article.Status),
		AuthorID:    article.AuthorID.String(),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}
	return e.producer.Publish(ctx, EventArticleUpdated, event)
}

// PublishArticleDeleted publishes an article deleted event
func (e *EventPublisher) PublishArticleDeleted(ctx context.Context, articleID string) error {
	event := models.ArticleDeletedEvent{
		ArticleID: articleID,
		DeletedAt: time.Now().Format(time.RFC3339),
	}
	return e.producer.Publish(ctx, EventArticleDeleted, event)
}

// PublishArticlePublished publishes an article published event
func (e *EventPublisher) PublishArticlePublished(ctx context.Context, article *models.Article) error {
	event := models.ArticlePublishedEvent{
		ArticleID:   article.ID.String(),
		Slug:        article.Slug,
		PublishedAt: time.Now().Format(time.RFC3339),
	}
	return e.producer.Publish(ctx, EventArticlePublished, event)
}

// PublishArticleArchived publishes an article archived event
func (e *EventPublisher) PublishArticleArchived(ctx context.Context, articleID string) error {
	event := models.ArticleArchivedEvent{
		ArticleID:  articleID,
		ArchivedAt: time.Now().Format(time.RFC3339),
	}
	return e.producer.Publish(ctx, EventArticleArchived, event)
}

// ========== REGION EVENTS ==========

// PublishRegionProvinceUpdated publishes a province updated event
func (e *EventPublisher) PublishRegionProvinceUpdated(ctx context.Context, province *models.Province) error {
	event := models.RegionProvinceUpdatedEvent{
		ProvinceID: province.ID,
		Name:       province.Name,
		UpdatedAt:  time.Now().Format(time.RFC3339),
	}
	return e.producer.Publish(ctx, EventRegionProvinceUpdated, event)
}

// PublishRegionRegencyUpdated publishes a regency updated event
func (e *EventPublisher) PublishRegionRegencyUpdated(ctx context.Context, regency *models.Regency) error {
	event := models.RegionRegencyUpdatedEvent{
		RegencyID:  regency.ID,
		ProvinceID: regency.ProvinceID,
		Name:       regency.Name,
		Type:       regency.Type,
		UpdatedAt:  time.Now().Format(time.RFC3339),
	}
	return e.producer.Publish(ctx, EventRegionRegencyUpdated, event)
}

// PublishRegionDistrictUpdated publishes a district updated event
func (e *EventPublisher) PublishRegionDistrictUpdated(ctx context.Context, district *models.District) error {
	event := models.RegionDistrictUpdatedEvent{
		DistrictID: district.ID,
		RegencyID:  district.RegencyID,
		Name:       district.Name,
		UpdatedAt:  time.Now().Format(time.RFC3339),
	}
	return e.producer.Publish(ctx, EventRegionDistrictUpdated, event)
}