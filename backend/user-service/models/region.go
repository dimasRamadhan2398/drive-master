package models

// Province represents a province from core-service
type Province struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// Regency represents a regency from core-service
type Regency struct {
	ID         uint   `json:"id"`
	ProvinceID uint   `json:"provinceId"`
	Name       string `json:"name"`
	Type       string `json:"type"`
}

// District represents a district from core-service
type District struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	RegencyID uint   `json:"regencyId"`
}