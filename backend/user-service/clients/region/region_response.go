package region

// APIResponse is the standard response format for region endpoints
type APIResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Data    RegionData `json:"data"`
}

type RegionData struct {
	Provinces []Province `json:"provinces,omitempty"`
	Regencies []Regency  `json:"regencies,omitempty"`
	Districts []District `json:"districts,omitempty"`
	Province  *Province  `json:"province,omitempty"`
	Regency   *Regency   `json:"regency,omitempty"`
}

// Master Data Models
type Province struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Regency struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ProvinceID string `json:"provinceId"`
	Type       string `json:"type"`
}

type District struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	RegencyID string `json:"regencyId"`
}

// List Response (jika API menggunakan list endpoint terpisah)
type ProvincesListResponse struct {
	Code    int        `json:"code"`
	Status  string     `json:"status"`
	Message string     `json:"message"`
	Data    []Province `json:"data"`
}

type RegenciesListResponse struct {
	Code    int       `json:"code"`
	Status  string    `json:"status"`
	Message string    `json:"message"`
	Data    []Regency `json:"data"`
}
