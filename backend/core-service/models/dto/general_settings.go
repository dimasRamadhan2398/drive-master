package dto

import (
	"core-service/models"
	"time"

	"github.com/google/uuid"
)

// ========== GENERAL SETTINGS DTOs ==========

// UpdateGeneralSettingsRequest represents the request to update general settings
type UpdateGeneralSettingsRequest struct {
	// Business Info
	BusinessName *string `json:"businessName"`
	Email        *string `json:"email"`
	Phone        *string `json:"phone"`
	Fax          *string `json:"fax"`
	WhatsApp     *string `json:"whatsApp"`
	Instagram    *string `json:"instagram"`
	Youtube      *string `json:"youtube"`
	MapDirection *string `json:"mapDirection"`

	// Address
	Address *string `json:"address"`

	// Operating Hours
	HoursMonFri     *string `json:"hoursMonFri"`
	HoursSatSun     *string `json:"hoursSatSun"`
	HoursNightShift *string `json:"hoursNightShift"`

	// Promo
	PromoEndDate *time.Time `json:"promoEndDate"`

	// Notification Settings
	NotifyEmail    *bool `json:"notifyEmail"`
	NotifySMS      *bool `json:"notifySms"`
	NotifyWhatsApp *bool `json:"notifyWhatsApp"`
}

// GeneralSettingsResponse represents the response for general settings
type GeneralSettingsResponse struct {
	ID          uuid.UUID `json:"id"`
	BusinessName string    `json:"businessName"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Fax         string    `json:"fax"`
	WhatsApp    string    `json:"whatsApp"`
	Instagram   string    `json:"instagram"`
	Youtube     string    `json:"youtube"`
	MapDirection string   `json:"mapDirection"`
	Address     string    `json:"address"`
	HoursMonFri string    `json:"hoursMonFri"`
	HoursSatSun string    `json:"hoursSatSun"`
	HoursNightShift string `json:"hoursNightShift"`
	PromoEndDate *time.Time `json:"promoEndDate"`
	NotifyEmail  bool    `json:"notifyEmail"`
	NotifySMS    bool    `json:"notifySms"`
	NotifyWhatsApp bool  `json:"notifyWhatsApp"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// ToResponse converts a GeneralSettings model to a response DTO
func ToGeneralSettingsResponse(settings *models.GeneralSettings) GeneralSettingsResponse {
	if settings == nil {
		return GeneralSettingsResponse{}
	}
	return GeneralSettingsResponse{
		ID:           settings.ID,
		BusinessName: settings.BusinessName,
		Email:        settings.Email,
		Phone:        settings.Phone,
		Fax:          settings.Fax,
		WhatsApp:     settings.WhatsApp,
		Instagram:    settings.Instagram,
		Youtube:      settings.Youtube,
		MapDirection: settings.MapDirection,
		Address:      settings.Address,
		HoursMonFri:  settings.HoursMonFri,
		HoursSatSun:  settings.HoursSatSun,
		HoursNightShift: settings.HoursNightShift,
		PromoEndDate: settings.PromoEndDate,
		NotifyEmail:  settings.NotifyEmail,
		NotifySMS:    settings.NotifySMS,
		NotifyWhatsApp: settings.NotifyWhatsApp,
		CreatedAt:    settings.CreatedAt,
		UpdatedAt:    settings.UpdatedAt,
	}
}