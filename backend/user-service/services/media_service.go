package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"path/filepath"

	apperrors "user-service/pkg/errors"

	"github.com/google/uuid"
	"github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/option"
	"github.com/imagekit-developer/imagekit-go/v2/packages/param"
)

type IMediaService interface {
	UploadMedia(ctx context.Context, input UploadMediaInput) (*MediaUploadResponse, error)
	UploadBase64Media(ctx context.Context, base64Data, fileName, folder string) (*MediaUploadResponse, error)
	UploadUserProfileImage(ctx context.Context, fileData []byte, userID uint, fileName string) (*MediaUploadResponse, error)
	DeleteMedia(ctx context.Context, fileID string) error
	GetMediaMetadata(ctx context.Context, fileID string) (*MediaUploadResponse, error)
	BulkDeleteMedia(ctx context.Context, fileIDs []string) error
}

type MediaService struct {
	imageKit imagekit.Client
}

func NewMediaService(privateKey string) IMediaService {
	client := imagekit.NewClient(
		option.WithPrivateKey(privateKey),
		option.WithPassword(privateKey),
	)
	return &MediaService{imageKit: client}
}

type MediaUploadResponse struct {
	URL         string `json:"url"`
	FileID      string `json:"fileId"`
	Name        string `json:"name"`
	Height      int    `json:"height,omitempty"`
	Width       int    `json:"width,omitempty"`
	FileType    string `json:"fileType"`
	CreatedAt   string `json:"createdAt"`
	Duration    float64 `json:"duration,omitempty"`   // For videos
	Size        int64  `json:"size,omitempty"`        // File size in bytes
	ThumbnailURL string `json:"thumbnailUrl,omitempty"` // Video thumbnail
	MIMEType    string `json:"mimeType,omitempty"`
}

type UploadMediaInput struct {
	FileData []byte
	FileName string
	Folder   string // e.g., "users/profiles", "lessons/videos"
}

// UploadMedia uploads a media file (image or video) to ImageKit and returns the response
func (s *MediaService) UploadMedia(ctx context.Context, input UploadMediaInput) (*MediaUploadResponse, error) {
	// Generate unique filename if not provided
	if input.FileName == "" {
		ext := detectFileExtension(input.FileData)
		input.FileName = fmt.Sprintf("%s%s", uuid.New().String(), ext)
	}

	// Create bytes.Reader as io.Reader for upload
	reader := bytes.NewReader(input.FileData)

	// Build params
	params := imagekit.FileUploadParams{
		File:     reader,
		FileName: input.FileName,
	}

	// Add folder if provided
	if input.Folder != "" {
		folder := input.Folder
		if len(folder) > 0 && folder[0] == '/' {
			folder = folder[1:] // Remove leading slash for API
		}
		params.Folder = param.NewOpt[string](folder)
	}

	resp, err := s.imageKit.Files.Upload(ctx, params)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}

	return &MediaUploadResponse{
		URL:      resp.URL,
		FileID:   resp.FileID,
		Name:     resp.Name,
		Height:   int(resp.Height),
		Width:    int(resp.Width),
		FileType: resp.FileType,
	}, nil
}

// UploadBase64Media uploads a base64-encoded media file to ImageKit
func (s *MediaService) UploadBase64Media(ctx context.Context, base64Data, fileName, folder string) (*MediaUploadResponse, error) {
	// Remove data URI prefix if present
	cleanedBase64, mimeType := extractBase64Data(base64Data)

	decoded, err := base64.StdEncoding.DecodeString(cleanedBase64)
	if err != nil {
		return nil, apperrors.ErrBadRequest
	}

	if fileName == "" {
		ext := mimeTypeToExtension(mimeType)
		fileName = fmt.Sprintf("%s%s", uuid.New().String(), ext)
	}

	resp, err := s.UploadMedia(ctx, UploadMediaInput{
		FileData: decoded,
		FileName: fileName,
		Folder:   folder,
	})
	if err != nil {
		return nil, err
	}

	resp.MIMEType = mimeType
	return resp, nil
}

// UploadUserProfileImage uploads a user profile image with proper naming
func (s *MediaService) UploadUserProfileImage(ctx context.Context, fileData []byte, userID uint, fileName string) (*MediaUploadResponse, error) {
	if fileName == "" {
		ext := detectFileExtension(fileData)
		fileName = fmt.Sprintf("user_%d_%s%s", userID, uuid.New().String()[:8], ext)
	}

	ext := filepath.Ext(fileName)
	if ext == "" {
		ext = ".jpg"
	}

	uniqueFileName := fmt.Sprintf("user_%d_%s%s", userID, uuid.New().String()[:8], ext)

	return s.UploadMedia(ctx, UploadMediaInput{
		FileData: fileData,
		FileName: uniqueFileName,
		Folder:   fmt.Sprintf("users/%d/profile", userID),
	})
}

// DeleteMedia deletes a media file from ImageKit by file ID
func (s *MediaService) DeleteMedia(ctx context.Context, fileID string) error {
	if fileID == "" {
		return apperrors.ErrBadRequest
	}
	err := s.imageKit.Files.Delete(ctx, fileID)
	if err != nil {
		return apperrors.ErrInternalServer
	}
	return nil
}

// GetMediaMetadata retrieves metadata for a media file by file ID
func (s *MediaService) GetMediaMetadata(ctx context.Context, fileID string) (*MediaUploadResponse, error) {
	if fileID == "" {
		return nil, apperrors.ErrBadRequest
	}

	resp, err := s.imageKit.Files.Get(ctx, fileID)
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}

	return &MediaUploadResponse{
		URL:      resp.URL,
		FileID:   resp.FileID,
		Name:     resp.Name,
		Height:   int(resp.Height),
		Width:    int(resp.Width),
		FileType: resp.FileType,
	}, nil
}

// BulkDeleteMedia deletes multiple media files by their file IDs
func (s *MediaService) BulkDeleteMedia(ctx context.Context, fileIDs []string) error {
	if len(fileIDs) == 0 {
		return apperrors.ErrBadRequest
	}

	params := imagekit.FileBulkDeleteParams{
		FileIDs: fileIDs,
	}
	_, err := s.imageKit.Files.Bulk.Delete(ctx, params)
	if err != nil {
		return apperrors.ErrInternalServer
	}
	return nil
}

// Helper functions

// detectFileExtension detects the file extension based on file content
func detectFileExtension(data []byte) string {
	if len(data) < 4 {
		return ".bin"
	}

	// PNG signature
	if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		return ".png"
	}
	// JPEG signature
	if data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return ".jpg"
	}
	// GIF signature
	if data[0] == 0x47 && data[1] == 0x49 && data[2] == 0x46 {
		return ".gif"
	}
	// WebP signature
	if data[0] == 0x52 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x46 &&
		len(data) >= 12 && data[8] == 0x57 && data[9] == 0x45 && data[10] == 0x42 && data[11] == 0x50 {
		return ".webp"
	}
	// MP4/MOV signature
	if len(data) >= 12 &&
		(data[4] == 0x66 && data[5] == 0x74 && data[6] == 0x79 && data[7] == 0x70 ||
			data[4] == 0x6D && data[5] == 0x6F && data[6] == 0x6F && data[7] == 0x76) {
		if data[8] == 0x6D && data[9] == 0x70 && data[10] == 0x34 && data[11] == 0x32 {
			return ".mp4"
		}
		if data[8] == 0x69 && data[9] == 0x73 && data[10] == 0x6F && data[11] == 0x6D {
			return ".mov"
		}
		if data[8] == 0x61 && data[9] == 0x76 && data[10] == 0x69 && data[11] == 0x66 {
			return ".avi"
		}
		return ".mp4"
	}
	// WebM signature
	if len(data) >= 4 && data[0] == 0x1A && data[1] == 0x45 && data[2] == 0xDF && data[3] == 0xA3 {
		return ".webm"
	}

	return ".bin"
}

// extractBase64Data removes data URI prefix and returns cleaned base64 and mime type
func extractBase64Data(base64Data string) (string, string) {
	mimeType := "application/octet-stream"

	// Check for various image MIME types
	if len(base64Data) > 22 && base64Data[:22] == "data:image/png;base64," {
		return base64Data[22:], "image/png"
	}
	if len(base64Data) > 23 && base64Data[:23] == "data:image/jpeg;base64," {
		return base64Data[23:], "image/jpeg"
	}
	if len(base64Data) > 19 && base64Data[:19] == "data:image/jpg;base64," {
		return base64Data[19:], "image/jpeg"
	}
	if len(base64Data) > 22 && base64Data[:22] == "data:image/webp;base64," {
		return base64Data[22:], "image/webp"
	}
	if len(base64Data) > 22 && base64Data[:22] == "data:image/gif;base64," {
		return base64Data[22:], "image/gif"
	}

	// Check for video MIME types
	if len(base64Data) > 23 && base64Data[:23] == "data:video/mp4;base64," {
		return base64Data[23:], "video/mp4"
	}
	if len(base64Data) > 25 && base64Data[:25] == "data:video/webm;base64," {
		return base64Data[25:], "video/webm"
	}
	if len(base64Data) > 25 && base64Data[:25] == "data:video/quicktime;base64," {
		return base64Data[25:], "video/quicktime"
	}
	if len(base64Data) > 24 && base64Data[:24] == "data:video/avi;base64," {
		return base64Data[24:], "video/avi"
	}

	// Check for PNG signature in base64 (iVBOR)
	if len(base64Data) > 5 && base64Data[:5] == "iVBOR" {
		return base64Data, "image/png"
	}

	// Check for JPEG signature in base64
	if len(base64Data) > 5 && base64Data[:5] == "/9j/4" {
		return base64Data, "image/jpeg"
	}

	return base64Data, mimeType
}

// mimeTypeToExtension converts MIME type to file extension
func mimeTypeToExtension(mimeType string) string {
	switch mimeType {
	case "image/png":
		return ".png"
	case "image/jpeg":
		return ".jpg"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	case "video/mp4":
		return ".mp4"
	case "video/webm":
		return ".webm"
	case "video/quicktime":
		return ".mov"
	case "video/avi":
		return ".avi"
	default:
		return ".bin"
	}
}
