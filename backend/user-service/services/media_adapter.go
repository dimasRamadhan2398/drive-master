package services

import core "core-service/services"

// Media adapter: re-export core-service media types so existing controllers and code
// in user-service can continue using services.IMediaService, services.UploadMediaInput, etc.

// Alias the interface and types from core-service
type IMediaService = core.IMediaService
type UploadMediaInput = core.UploadMediaInput
type MediaUploadResponse = core.MediaUploadResponse

// Optionally, you can add helper wrapper functions here in the future if
// translation between models is needed.
