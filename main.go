package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rivo/tview"
)

type GeocodeResponse struct {
	Results []struct {
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Location struct {
				Lat float64
				Lng float64
			}
		}
	}
}

type WeatherPointResponse struct {
	Properties struct {
		Forecast string `json:"forecast"`
	}
}

type ForecastResponse struct {
	Properties struct {
		Periods []struct {
			Number           int    `json:"number"`
			Name             string `json:"name"`
			StartTime        string `json:"startTime"`
			EndTime          string `json:"endTime"`
			IsDaytime        bool   `json:"isDaytime"`
			Temperature      int    `json:"temperature"`
			TemperatureUnit  string `json:"temperatureUnit"`
			WindSpeed        string `json:"windSpeed"`
			WindDirection    string `json:"windDirection"`
			Icon             string `json:"icon"`
			ShortForecast    string `json:"shortForecast"`
			DetailedForecast string `json:"detailedForecast"`
		}
	}
}

func main() {
	app := tview.NewApplication()
	citySelectForm := tview.NewForm().
		AddInputField("City Name", "Oxford", 20, nil, nil).
		AddInputField("State Name", "OH", 20, nil, nil).
		AddButton("Submit", func() {
			app.Stop()
		})

	citySelectForm.SetBorder(true).SetTitle("Select a City and State").SetTitleAlign(tview.AlignLeft)

	if err := app.SetRoot(citySelectForm, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	fmt.Printf("%s %s\n", citySelectForm.GetFormItem(0).(*tview.InputField).GetText(), citySelectForm.GetFormItem(1).(*tview.InputField).GetText())

	//apiTest := googleQuery(citySelectForm.GetFormItem(0).(*tview.InputField).GetText(), citySelectForm.GetFormItem(1).(*tview.InputField).GetText())
	//fmt.Printf("%+v\n", apiTest.Results[0].FormattedAddress)
	//fmt.Printf("%+v\n", apiTest.Results[0].Geometry.Location.Lat)
	//fmt.Printf("%+v\n", apiTest.Results[0].Geometry.Location.Lng)

	//latCnv := strconv.FormatFloat(apiTest.Results[0].Geometry.Location.Lat, 'f', 5, 64)
	//lngCnv := strconv.FormatFloat(apiTest.Results[0].Geometry.Location.Lng, 'f', 5, 64)
	//weatherApiTest := weatherPointQuery(latCnv, lngCnv)
	//fmt.Printf("%v\n", weatherApiTest.Properties.Forecast)

	//forecastTest := forecastQuery(weatherApiTest.Properties.Forecast)
	//fmt.Println(forecastTest)
}

func googleQuery(city string, state string) GeocodeResponse {
	googleApiEndpoint := "https://maps.googleapis.com/maps/api/geocode/json?address="
	googleApiKey := "&key=CHANGEME"

	client := &http.Client{}
	reqUri := googleApiEndpoint + city + "+" + state + googleApiKey
	request, err := http.NewRequest("GET", reqUri, nil)
	if err != nil {
		fmt.Print(err.Error())
	}

	response, err := client.Do(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer response.Body.Close()

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	var responseObj GeocodeResponse
	json.Unmarshal(responseBytes, &responseObj)

	return responseObj
}

func weatherPointQuery(lat string, lng string) WeatherPointResponse {
	weatherPtApiEndpoint := "https://api.weather.gov/points/"

	client := &http.Client{}
	reqUri := weatherPtApiEndpoint + lat + "," + lng

	request, err := http.NewRequest("GET", reqUri, nil)
	if err != nil {
		fmt.Print(err.Error())
	}

	response, err := client.Do(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer response.Body.Close()

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	var responseObj WeatherPointResponse
	json.Unmarshal(responseBytes, &responseObj)

	return responseObj
}

func forecastQuery(url string) ForecastResponse {
	forecastApiEndpoint := url

	client := &http.Client{}
	reqUri := forecastApiEndpoint

	request, err := http.NewRequest("GET", reqUri, nil)
	if err != nil {
		fmt.Print(err.Error())
	}

	response, err := client.Do(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer response.Body.Close()

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	var responseObj ForecastResponse
	json.Unmarshal(responseBytes, &responseObj)

	return responseObj
}
