package models

type OpenMeteoAPI struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Elevation float64 `json:"elevation"`
	Hourly    Hourly  `json:"hourly"`
}

type Hourly struct {
	Time              []int64   `json:"time"`
	Temperature2m     []float64 `json:"temperature_2m"`
	PrecipProbability []int     `json:"precipitation_probability"`
	Precipitation     []float64 `json:"precipitation"`
	Visibility        []float64 `json:"visibility"`
	WindSpeed10m      []float64 `json:"wind_speed_10m"`
}
