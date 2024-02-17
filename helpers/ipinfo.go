package helpers

import (
	"encoding/json"

	"github.com/g-bolmida/weather-app/models"
)

func IPInfo() models.IPInfoAPI {
	body, err := APIRequest("https://ipinfo.io/json")
	if err != nil {
		panic(err)
	}

	var apiResponse models.IPInfoAPI
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		panic(err)
	}

	return apiResponse
}
