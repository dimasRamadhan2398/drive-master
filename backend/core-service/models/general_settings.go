package models

import (
	"time"

	"github.com/google/uuid"
)

// GeneralSettings holds all general business information and configuration
type GeneralSettings struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`

	// Business Info
	BusinessName string `json:"businessName" gorm:"size:255"`
	Email        string `json:"email" gorm:"size:255"`
	Phone        string `json:"phone" gorm:"size:50"`
	Fax          string `json:"fax" gorm:"size:50"`
	WhatsApp     string `json:"whatsApp" gorm:"size:50"`

	// Address
	Address string `json:"address" gorm:"type:text"`

	// Operating Hours (stored as strings for flexibility - e.g., "08:00 - 17:00")
	HoursMonFri     string `json:"hoursMonFri" gorm:"size:100"`      // Monday - Friday
	HoursSatSun     string `json:"hoursSatSun" gorm:"size:100"`      // Saturday - Sunday
	HoursNightShift string `json:"hoursNightShift" gorm:"size:100"`  // Night shift hours

	// Promo
	PromoEndDate *time.Time `json:"promoEndDate" gorm:"type:date"`

	// Notification Settings
	NotifyEmail    bool `json:"notifyEmail" gorm:"default:true"`
	NotifySMS      bool `json:"notifySms" gorm:"default:false"`
	NotifyWhatsApp bool `json:"notifyWhatsApp" gorm:"default:false"`

	// Timestamps
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}