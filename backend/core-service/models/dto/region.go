package dto

type RegionResponse struct {
	Code    int       `json:"code"`
	Status  string    `json:"status"`
	Message string    `json:"message"`
	Data    RegionData `json:"data"`
}

type RegionData struct {
    Provinces  []Province  `json:"provinces,omitempty"`
    Kabupatens []Kabupaten `json:"kabupatens,omitempty"`
    Kecamatans []Kecamatan `json:"kecamatans,omitempty"`
    Province   *Province   `json:"province,omitempty"`
    Kabupaten  *Kabupaten  `json:"kabupaten,omitempty"`
    Kecamatan  *Kecamatan  `json:"kecamatan,omitempty"`
}

// Master Data Models
type Province struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

type Kabupaten struct {
    ID         string `json:"id"`
    Name       string `json:"name"`
    ProvinceID string `json:"province_id"`
}

type Kecamatan struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    KabupatenID string `json:"kabupaten_id"`
}

// List Response (jika API menggunakan list endpoint terpisah)
type ProvincesListResponse struct {
    Code    int        `json:"code"`
    Status  string     `json:"status"`
    Message string     `json:"message"`
    Data    []Province `json:"data"`
}

type KabupatensListResponse struct {
    Code    int         `json:"code"`
    Status  string      `json:"status"`
    Message string      `json:"message"`
    Data    []Kabupaten `json:"data"`
}