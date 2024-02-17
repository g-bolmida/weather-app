package models

type IPInfoAPI struct {
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Postal   string `json:"postal"`
	Timezone string `json:"timezone"`
}
