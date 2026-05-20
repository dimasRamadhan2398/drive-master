package dto

import "slices"

// MediaType represents the type of media file
type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
	MediaTypeFile  MediaType = "file"
)

// MediaUploadResponse represents the response after uploading a media file
type MediaUploadResponse struct {
	URL          string     `json:"url"`
	FileID       string     `json:"fileId"`
	Name         string     `json:"name"`
	Height       int        `json:"height,omitempty"`
	Width        int        `json:"width,omitempty"`
	FileType     string     `json:"fileType"`
	CreatedAt    string     `json:"createdAt"`
	Duration     float64    `json:"duration,omitempty"`    // For videos
	Size         int64      `json:"size,omitempty"`        // File size in bytes
	ThumbnailURL string     `json:"thumbnailUrl,omitempty"` // Video thumbnail
	MIMEType     string     `json:"mimeType,omitempty"`
	MediaType    MediaType  `json:"mediaType,omitempty"`   // Derived from MIME type
}

// UploadMediaRequest represents the request to upload a media file
type UploadMediaRequest struct {
	FileData []byte `json:"fileData" validate:"required"`
	FileName string `json:"fileName"`
	Folder   string `json:"folder"`
}

// UploadBase64MediaRequest represents the request to upload a base64-encoded media file
type UploadBase64MediaRequest struct {
	Base64Data string `json:"base64Data" validate:"required"`
	FileName   string `json:"fileName"`
	Folder     string `json:"folder"`
}

// DeleteMediaRequest represents the request to delete a media file
type DeleteMediaRequest struct {
	FileID string `json:"fileId" validate:"required"`
}

// BulkDeleteMediaRequest represents the request to delete multiple media files
type BulkDeleteMediaRequest struct {
	FileIDs []string `json:"fileIds" validate:"required,min=1"`
}

// GetMediaMetadataRequest represents the request to get media metadata
type GetMediaMetadataRequest struct {
	FileID string `json:"fileId" validate:"required"`
}

// MediaListResponse represents a paginated list of media files
type MediaListResponse struct {
	Items      []MediaUploadResponse `json:"items"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"pageSize"`
	TotalPages int                   `json:"totalPages"`
}

// SupportedMIMETypes lists all supported MIME types for media uploads
var SupportedMIMETypes = []string{
	// Images
	"image/png",
	"image/jpeg",
	"image/jpg",
	"image/gif",
	"image/webp",
	"image/svg+xml",
	// Videos
	"video/mp4",
	"video/webm",
	"video/quicktime",
	"video/x-msvideo",
	"video/x-matroska",
}

// IsImageMIMEType checks if the given MIME type is an image
func IsImageMIMEType(mimeType string) bool {
	imageTypes := []string{
		"image/png",
		"image/jpeg",
		"image/jpg",
		"image/gif",
		"image/webp",
		"image/svg+xml",
	}
	return slices.Contains(imageTypes, mimeType)
}

// IsVideoMIMEType checks if the given MIME type is a video
func IsVideoMIMEType(mimeType string) bool {
	videoTypes := []string{
		"video/mp4",
		"video/webm",
		"video/quicktime",
		"video/x-msvideo",
		"video/x-matroska",
	}
	return slices.Contains(videoTypes, mimeType)
}

// GetMediaTypeFromMIME returns the MediaType based on MIME type
func GetMediaTypeFromMIME(mimeType string) MediaType {
	if IsImageMIMEType(mimeType) {
		return MediaTypeImage
	}
	if IsVideoMIMEType(mimeType) {
		return MediaTypeVideo
	}
	return MediaTypeFile
}

// Legacy types for backward compatibility
type ImageUploadResponse = MediaUploadResponse
type UploadImageInput = UploadMediaRequest
