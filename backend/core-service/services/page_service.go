package services

import (
	"context"
	"core-service/models"
	"core-service/models/dto"
	"core-service/repositories"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
)

type IPageService interface {
	GetAllPages(ctx context.Context) ([]dto.PageResponse, error)
	GetPageByID(ctx context.Context, id uuid.UUID) (*dto.PageResponse, error)
	GetPageBySlug(ctx context.Context, slug string) (*dto.PageResponse, error)
	CreatePage(ctx context.Context, req *dto.CreatePageRequest) (*dto.PageResponse, error)
	UpdatePage(ctx context.Context, id uuid.UUID, req *dto.UpdatePageRequest) (*dto.PageResponse, error)
	DeletePage(ctx context.Context, id uuid.UUID) error
}

type PageService struct {
	pageRepo     repositories.IPageRepository
	mediaService IMediaService
}

func NewPageService(pageRepo repositories.IPageRepository, mediaService IMediaService) IPageService {
	return &PageService{
		pageRepo:     pageRepo,
		mediaService: mediaService,
	}
}

func (s *PageService) toResponse(page *models.Page) dto.PageResponse {
	return dto.PageResponse{
		ID:          page.ID,
		Title:       page.Title,
		Slug:        page.Slug,
		Status:      string(page.Status),
		Sections:    page.Sections,
		LastUpdated: page.UpdatedAt.Format("Jan _2, 2006"),
		CreatedAt:   page.CreatedAt,
		UpdatedAt:   page.UpdatedAt,
	}
}

func (s *PageService) GetAllPages(ctx context.Context) ([]dto.PageResponse, error) {
	pages, err := s.pageRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]dto.PageResponse, len(pages))
	for i, page := range pages {
		response[i] = s.toResponse(&page)
	}
	return response, nil
}

func (s *PageService) GetPageByID(ctx context.Context, id uuid.UUID) (*dto.PageResponse, error) {
	page, err := s.pageRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if page == nil {
		return nil, nil
	}
	res := s.toResponse(page)
	return &res, nil
}

func (s *PageService) GetPageBySlug(ctx context.Context, slug string) (*dto.PageResponse, error) {
	page, err := s.pageRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if page == nil {
		return nil, nil
	}
	res := s.toResponse(page)
	return &res, nil
}

func (s *PageService) CreatePage(ctx context.Context, req *dto.CreatePageRequest) (*dto.PageResponse, error) {
	page := &models.Page{
		Title:    req.Title,
		Slug:     req.Slug,
		Status:   models.PageStatus(req.Status),
		Sections: "[]", // Default to empty array
	}

	if err := s.pageRepo.Create(ctx, page); err != nil {
		return nil, err
	}

	res := s.toResponse(page)
	return &res, nil
}

func (s *PageService) uploadLocalImages(ctx context.Context, sectionsJSON string) (string, error) {
	if sectionsJSON == "" || sectionsJSON == "[]" {
		return sectionsJSON, nil
	}

	var sections []map[string]interface{}
	if err := json.Unmarshal([]byte(sectionsJSON), &sections); err != nil {
		return sectionsJSON, err
	}

	modified := false
	for _, section := range sections {
		sectionType, _ := section["type"].(string)
		data, ok := section["data"].(map[string]interface{})
		if !ok {
			continue
		}

		if sectionType == "hero" {
			if bgImage, exists := data["bgImage"].(string); exists && strings.HasPrefix(bgImage, "data:image/") {
				resp, err := s.mediaService.UploadBase64Media(ctx, bgImage, "", "pages")
				if err == nil && resp != nil {
					data["bgImage"] = resp.URL
					modified = true
				}
			}
		} else if sectionType == "image_text" {
			if image, exists := data["image"].(string); exists && strings.HasPrefix(image, "data:image/") {
				resp, err := s.mediaService.UploadBase64Media(ctx, image, "", "pages")
				if err == nil && resp != nil {
					data["image"] = resp.URL
					modified = true
				}
			}
		}
	}

	if !modified {
		return sectionsJSON, nil
	}

	updatedBytes, err := json.Marshal(sections)
	if err != nil {
		return sectionsJSON, err
	}

	return string(updatedBytes), nil
}

func (s *PageService) UpdatePage(ctx context.Context, id uuid.UUID, req *dto.UpdatePageRequest) (*dto.PageResponse, error) {
	page, err := s.pageRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if page == nil {
		return nil, nil
	}

	if req.Title != nil {
		page.Title = *req.Title
	}
	if req.Slug != nil {
		page.Slug = *req.Slug
	}
	if req.Status != nil {
		page.Status = models.PageStatus(*req.Status)
	}
	if req.Sections != nil {
		uploadedSections, err := s.uploadLocalImages(ctx, *req.Sections)
		if err == nil {
			page.Sections = uploadedSections
		} else {
			page.Sections = *req.Sections
		}
	}

	page.UpdatedAt = time.Now()

	if err := s.pageRepo.Update(ctx, page); err != nil {
		return nil, err
	}

	res := s.toResponse(page)
	return &res, nil
}

func (s *PageService) DeletePage(ctx context.Context, id uuid.UUID) error {
	return s.pageRepo.Delete(ctx, id)
}
