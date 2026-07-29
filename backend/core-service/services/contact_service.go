package services

import (
	"context"
	"core-service/models"
	"core-service/models/dto"
	"core-service/repositories"
)

type IContactService interface {
	CreateInquiry(ctx context.Context, req *dto.CreateContactInquiryRequest) (*models.ContactInquiry, error)
	GetAllInquiries(ctx context.Context) ([]models.ContactInquiry, error)
}

type ContactService struct {
	contactRepo repositories.IContactRepository
}

func NewContactService(contactRepo repositories.IContactRepository) IContactService {
	return &ContactService{
		contactRepo: contactRepo,
	}
}

func (s *ContactService) CreateInquiry(ctx context.Context, req *dto.CreateContactInquiryRequest) (*models.ContactInquiry, error) {
	contact := &models.ContactInquiry{
		Name:    req.Name,
		Email:   req.Email,
		Subject: req.Subject,
		Message: req.Message,
	}

	if err := s.contactRepo.Create(ctx, contact); err != nil {
		return nil, err
	}

	return contact, nil
}

func (s *ContactService) GetAllInquiries(ctx context.Context) ([]models.ContactInquiry, error) {
	return s.contactRepo.FindAll(ctx)
}
