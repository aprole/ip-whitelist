package models

type HTTPRequest struct {
	IP               string   `json:"ip"`
	AllowedCountries []string `json:"allowedCountries"`
}

type HTTPResponse struct {
	Accepted    bool   `json:"accepted"`
	IP          string `json:"ip"`
	Country     string `json:"country"`
	CountryName string `json:"countryName"`
}
