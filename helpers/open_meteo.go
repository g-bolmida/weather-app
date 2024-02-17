package helpers

import (
	"encoding/json"
	"net/url"

	"github.com/g-bolmida/weather-app/models"
)

func GetWeather(lat string, long string, timezone string) models.OpenMeteoAPI {
	url := "https://api.open-meteo.com/v1/forecast?latitude=" + lat + "&longitude=" + long + "&hourly=temperature_2m,precipitation_probability,precipitation,visibility,wind_speed_10m&temperature_unit=fahrenheit&wind_speed_unit=mph&precipitation_unit=inch&timeformat=unixtime&timezone=" + url.QueryEscape(timezone) + "&forecast_days=1"

	body, err := APIRequest(url)
	if err != nil {
		panic(err)
	}

	var apiResponse models.OpenMeteoAPI
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		panic(err)
	}

	return apiResponse
}
