package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"core-service/models"
	"core-service/repositories"
)

type IEventService interface {
	// User events (consumed from user-service)
	HandleUserCreated(ctx context.Context, event models.UserCreatedEvent) error
	HandleUserUpdated(ctx context.Context, event models.UserUpdatedEvent) error
	HandleUserDeleted(ctx context.Context, event models.UserDeletedEvent) error

	// Car events (consumed by other services)
	HandleCarCreated(ctx context.Context, event models.CarCreatedEvent) error
	HandleCarUpdated(ctx context.Context, event models.CarUpdatedEvent) error
	HandleCarDeleted(ctx context.Context, event models.CarDeletedEvent) error

	// Package events (consumed by other services)
	HandlePackageCreated(ctx context.Context, event models.PackageCreatedEvent) error
	HandlePackageUpdated(ctx context.Context, event models.PackageUpdatedEvent) error
	HandlePackageDeleted(ctx context.Context, event models.PackageDeletedEvent) error

	// Article events (consumed by other services like search, notification)
	HandleArticleCreated(ctx context.Context, event models.ArticleCreatedEvent) error
	HandleArticleUpdated(ctx context.Context, event models.ArticleUpdatedEvent) error
	HandleArticleDeleted(ctx context.Context, event models.ArticleDeletedEvent) error
	HandleArticlePublished(ctx context.Context, event models.ArticlePublishedEvent) error
	HandleArticleArchived(ctx context.Context, event models.ArticleArchivedEvent) error

	// Region events (consumed by other services)
	HandleRegionProvinceUpdated(ctx context.Context, event models.RegionProvinceUpdatedEvent) error
	HandleRegionRegencyUpdated(ctx context.Context, event models.RegionRegencyUpdatedEvent) error
	HandleRegionDistrictUpdated(ctx context.Context, event models.RegionDistrictUpdatedEvent) error

	// System events (consumed by caching, monitoring)
	HandleCacheInvalidated(ctx context.Context, event models.CacheInvalidatedEvent) error
}

type EventService struct {
	eventRepo repositories.IEventRepository
	cacheRepo repositories.ICacheRepository
}

func NewEventService(eventRepo repositories.IEventRepository, cacheRepo repositories.ICacheRepository) IEventService {
	return &EventService{
		eventRepo: eventRepo,
		cacheRepo: cacheRepo,
	}
}

// storeProcessedEvent stores an event in the database for audit/traceability
func (s *EventService) storeProcessedEvent(ctx context.Context, eventType string, payload interface{}) error {
	if s.eventRepo == nil {
		return nil
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal event payload: %w", err)
	}

	return s.eventRepo.SaveProcessedEvent(ctx, eventType, string(payloadBytes))
}

// HandleUserCreated handles user.created events
func (s *EventService) HandleUserCreated(ctx context.Context, event models.UserCreatedEvent) error {
	log.Printf("[EventService] Processing user.created for userId=%d, email=%s", event.UserID, event.Email)

	// TODO: Add business logic based on what other services need when a user is created
	// Examples:
	// - Invalidate user-related caches
	// - Send welcome notification
	// - Create initial user preferences
	// - Sync user data to search service

	if err := s.storeProcessedEvent(ctx, "user.created", event); err != nil {
		log.Printf("[EventService] Failed to store user.created event: %v", err)
		// Don't fail the handler if storing fails
	}

	return nil
}

// HandleUserUpdated handles user.updated events
func (s *EventService) HandleUserUpdated(ctx context.Context, event models.UserUpdatedEvent) error {
	log.Printf("[EventService] Processing user.updated for userId=%d", event.UserID)

	// Invalidate user cache if available
	if s.cacheRepo != nil {
		cacheKey := fmt.Sprintf("user:%d", event.UserID)
		if err := s.cacheRepo.Delete(ctx, cacheKey); err != nil {
			log.Printf("[EventService] Failed to invalidate user cache: %v", err)
		}
	}

	return s.storeProcessedEvent(ctx, "user.updated", event)
}

// HandleUserDeleted handles user.deleted events
func (s *EventService) HandleUserDeleted(ctx context.Context, event models.UserDeletedEvent) error {
	log.Printf("[EventService] Processing user.deleted for userId=%d", event.UserID)

	// Invalidate all related caches
	if s.cacheRepo != nil {
		// Pattern-based invalidation for user's related caches
		patterns := []string{
			fmt.Sprintf("user:%d", event.UserID),
			fmt.Sprintf("user:%d:profile", event.UserID),
			fmt.Sprintf("user:%d:preferences", event.UserID),
		}
		for _, pattern := range patterns {
			if err := s.cacheRepo.Delete(ctx, pattern); err != nil {
				log.Printf("[EventService] Failed to invalidate cache %s: %v", pattern, err)
			}
		}
	}

	return s.storeProcessedEvent(ctx, "user.deleted", event)
}

// HandleCarCreated handles car.created events
func (s *EventService) HandleCarCreated(ctx context.Context, event models.CarCreatedEvent) error {
	log.Printf("[EventService] Processing car.created for carId=%s", event.CarID)

	// Notify search service to index the new car
	// Notify notification service for admin alerts if needed

	return s.storeProcessedEvent(ctx, "car.created", event)
}

// HandleCarUpdated handles car.updated events
func (s *EventService) HandleCarUpdated(ctx context.Context, event models.CarUpdatedEvent) error {
	log.Printf("[EventService] Processing car.updated for carId=%s", event.CarID)

	// Invalidate car-related caches
	if s.cacheRepo != nil {
		cacheKey := fmt.Sprintf("car:%s", event.CarID)
		s.cacheRepo.Delete(ctx, cacheKey)
	}

	return s.storeProcessedEvent(ctx, "car.updated", event)
}

// HandleCarDeleted handles car.deleted events
func (s *EventService) HandleCarDeleted(ctx context.Context, event models.CarDeletedEvent) error {
	log.Printf("[EventService] Processing car.deleted for carId=%s", event.CarID)

	// Invalidate car-related caches
	if s.cacheRepo != nil {
		cacheKey := fmt.Sprintf("car:%s", event.CarID)
		s.cacheRepo.Delete(ctx, cacheKey)
	}

	return s.storeProcessedEvent(ctx, "car.deleted", event)
}

// HandlePackageCreated handles package.created events
func (s *EventService) HandlePackageCreated(ctx context.Context, event models.PackageCreatedEvent) error {
	log.Printf("[EventService] Processing package.created for packageId=%s", event.PackageID)

	// Notify search service to index the new package

	return s.storeProcessedEvent(ctx, "package.created", event)
}

// HandlePackageUpdated handles package.updated events
func (s *EventService) HandlePackageUpdated(ctx context.Context, event models.PackageUpdatedEvent) error {
	log.Printf("[EventService] Processing package.updated for packageId=%s", event.PackageID)

	if s.cacheRepo != nil {
		cacheKey := fmt.Sprintf("package:%s", event.PackageID)
		s.cacheRepo.Delete(ctx, cacheKey)
	}

	return s.storeProcessedEvent(ctx, "package.updated", event)
}

// HandlePackageDeleted handles package.deleted events
func (s *EventService) HandlePackageDeleted(ctx context.Context, event models.PackageDeletedEvent) error {
	log.Printf("[EventService] Processing package.deleted for packageId=%s", event.PackageID)

	if s.cacheRepo != nil {
		cacheKey := fmt.Sprintf("package:%s", event.PackageID)
		s.cacheRepo.Delete(ctx, cacheKey)
	}

	return s.storeProcessedEvent(ctx, "package.deleted", event)
}

// HandleArticleCreated handles article.created events
func (s *EventService) HandleArticleCreated(ctx context.Context, event models.ArticleCreatedEvent) error {
	log.Printf("[EventService] Processing article.created for articleId=%s", event.ArticleID)

	// Notify search service to index the new article
	// Notify notification service for new content alerts

	return s.storeProcessedEvent(ctx, "article.created", event)
}

// HandleArticleUpdated handles article.updated events
func (s *EventService) HandleArticleUpdated(ctx context.Context, event models.ArticleUpdatedEvent) error {
	log.Printf("[EventService] Processing article.updated for articleId=%s", event.ArticleID)

	if s.cacheRepo != nil {
		cacheKey := fmt.Sprintf("article:%s", event.ArticleID)
		s.cacheRepo.Delete(ctx, cacheKey)
	}

	return s.storeProcessedEvent(ctx, "article.updated", event)
}

// HandleArticleDeleted handles article.deleted events
func (s *EventService) HandleArticleDeleted(ctx context.Context, event models.ArticleDeletedEvent) error {
	log.Printf("[EventService] Processing article.deleted for articleId=%s", event.ArticleID)

	if s.cacheRepo != nil {
		cacheKey := fmt.Sprintf("article:%s", event.ArticleID)
		s.cacheRepo.Delete(ctx, cacheKey)
	}

	return s.storeProcessedEvent(ctx, "article.deleted", event)
}

// HandleArticlePublished handles article.published events
func (s *EventService) HandleArticlePublished(ctx context.Context, event models.ArticlePublishedEvent) error {
	log.Printf("[EventService] Processing article.published for articleId=%s", event.ArticleID)

	// Notify notification service to broadcast new article
	// Notify social media integration if auto-post enabled

	return s.storeProcessedEvent(ctx, "article.published", event)
}

// HandleArticleArchived handles article.archived events
func (s *EventService) HandleArticleArchived(ctx context.Context, event models.ArticleArchivedEvent) error {
	log.Printf("[EventService] Processing article.archived for articleId=%s", event.ArticleID)

	if s.cacheRepo != nil {
		cacheKey := fmt.Sprintf("article:%s", event.ArticleID)
		s.cacheRepo.Delete(ctx, cacheKey)
	}

	return s.storeProcessedEvent(ctx, "article.archived", event)
}

// HandleRegionProvinceUpdated handles region.province.updated events
func (s *EventService) HandleRegionProvinceUpdated(ctx context.Context, event models.RegionProvinceUpdatedEvent) error {
	log.Printf("[EventService] Processing region.province.updated for provinceId=%d", event.ProvinceID)

	if s.cacheRepo != nil {
		cacheKey := fmt.Sprintf("region:province:%d", event.ProvinceID)
		s.cacheRepo.Delete(ctx, cacheKey)
	}

	return s.storeProcessedEvent(ctx, "region.province.updated", event)
}

// HandleRegionRegencyUpdated handles region.regency.updated events
func (s *EventService) HandleRegionRegencyUpdated(ctx context.Context, event models.RegionRegencyUpdatedEvent) error {
	log.Printf("[EventService] Processing region.regency.updated for regencyId=%d", event.RegencyID)

	if s.cacheRepo != nil {
		cacheKey := fmt.Sprintf("region:regency:%d", event.RegencyID)
		s.cacheRepo.Delete(ctx, cacheKey)
	}

	return s.storeProcessedEvent(ctx, "region.regency.updated", event)
}

// HandleRegionDistrictUpdated handles region.district.updated events
func (s *EventService) HandleRegionDistrictUpdated(ctx context.Context, event models.RegionDistrictUpdatedEvent) error {
	log.Printf("[EventService] Processing region.district.updated for districtId=%d", event.DistrictID)

	if s.cacheRepo != nil {
		cacheKey := fmt.Sprintf("region:district:%d", event.DistrictID)
		s.cacheRepo.Delete(ctx, cacheKey)
	}

	return s.storeProcessedEvent(ctx, "region.district.updated", event)
}

// HandleCacheInvalidated handles cache.invalidated events
func (s *EventService) HandleCacheInvalidated(ctx context.Context, event models.CacheInvalidatedEvent) error {
	log.Printf("[EventService] Processing cache.invalidated for key=%s", event.CacheKey)

	// Log for monitoring/analytics purposes

	return s.storeProcessedEvent(ctx, "cache.invalidated", event)
}
