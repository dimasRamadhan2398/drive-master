package services

import (
	"context"
	"encoding/json"
	"fmt"

	"core-service/models"
	"core-service/repositories"
)

type EventService struct {
	eventRepo *repositories.EventRepository
	cacheRepo *repositories.CacheRepository
}

func NewEventService(eventRepo *repositories.EventRepository, cacheRepo *repositories.CacheRepository) *EventService {
	return &EventService{
		eventRepo: eventRepo,
		cacheRepo: cacheRepo,
	}
}

func (s *EventService) HandleUserCreated(ctx context.Context, event models.UserCreatedEvent) error {
	payloadBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	if err := s.eventRepo.SaveProcessedEvent("user.created", string(payloadBytes)); err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("user:%d", event.UserID)
	return s.cacheRepo.SetUserProfile(ctx, cacheKey, string(payloadBytes))
}
