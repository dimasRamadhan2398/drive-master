package services

import (
	"context"
	"core-service/models"
	"core-service/repositories"

	"github.com/google/uuid"
)

// IFAQService defines the interface for FAQ service
type IFAQService interface {
	GetAllFAQs(ctx context.Context) ([]models.FAQ, error)
	GetActiveFAQs(ctx context.Context) ([]models.FAQ, error)
	GetFAQByID(ctx context.Context, id uuid.UUID) (*models.FAQ, error)
	CreateFAQ(ctx context.Context, question, answer, category string, order int) (*models.FAQ, error)
	UpdateFAQ(ctx context.Context, id uuid.UUID, question, answer, category string, order int, isActive bool) (*models.FAQ, error)
	DeleteFAQ(ctx context.Context, id uuid.UUID) error
}

// FAQService implements IFAQService
type FAQService struct {
	faqRepo repositories.IFAQRepository
}

// NewFAQService creates a new FAQ service
func NewFAQService(faqRepo repositories.IFAQRepository) IFAQService {
	return &FAQService{
		faqRepo: faqRepo,
	}
}

// GetAllFAQs retrieves all FAQs (including inactive)
func (s *FAQService) GetAllFAQs(ctx context.Context) ([]models.FAQ, error) {
	return s.faqRepo.FindAll(ctx)
}

// GetActiveFAQs retrieves all active FAQs
func (s *FAQService) GetActiveFAQs(ctx context.Context) ([]models.FAQ, error) {
	return s.faqRepo.FindActive(ctx)
}

// GetFAQByID retrieves an FAQ by ID
func (s *FAQService) GetFAQByID(ctx context.Context, id uuid.UUID) (*models.FAQ, error) {
	return s.faqRepo.FindByID(ctx, id)
}

// CreateFAQ creates a new FAQ
func (s *FAQService) CreateFAQ(ctx context.Context, question, answer, category string, order int) (*models.FAQ, error) {
	faq := &models.FAQ{
		Question: question,
		Answer:   answer,
		Category: category,
		Order:    order,
		IsActive: true,
	}

	if err := s.faqRepo.Create(ctx, faq); err != nil {
		return nil, err
	}

	return faq, nil
}

// UpdateFAQ updates an existing FAQ
func (s *FAQService) UpdateFAQ(ctx context.Context, id uuid.UUID, question, answer, category string, order int, isActive bool) (*models.FAQ, error) {
	faq, err := s.faqRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if faq == nil {
		return nil, nil
	}

	faq.Question = question
	faq.Answer = answer
	faq.Category = category
	faq.Order = order
	faq.IsActive = isActive

	if err := s.faqRepo.Update(ctx, faq); err != nil {
		return nil, err
	}

	return faq, nil
}

// DeleteFAQ soft-deletes an FAQ
func (s *FAQService) DeleteFAQ(ctx context.Context, id uuid.UUID) error {
	faq, err := s.faqRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if faq == nil {
		return nil
	}

	return s.faqRepo.DeleteSoft(ctx, faq)
}