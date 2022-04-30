package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gdamore/tcell/v2"
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

	weatherAppUI := tview.NewApplication()

	var citySelection *tview.Form
	var columns *tview.Table
	var flex *tview.Flex

	columns = tview.NewTable().SetBorders(true)
	columns.SetBorder(true).SetTitle("Select A City and State")

	citySelection = tview.NewForm().
		AddInputField("City Name", "", 20, nil, nil).
		AddInputField("State Name", "", 20, nil, nil).
		AddButton("Submit", func() {
			cityVal := citySelection.GetFormItem(0).(*tview.InputField).GetText()
			stateVal := citySelection.GetFormItem(1).(*tview.InputField).GetText()
			
			if cityVal == "" || stateVal == "" {
				log.Fatal("City or State not entered...\n")
				weatherAppUI.Stop()
			}
			
			getCoords, getForecast := apiCalls(cityVal, stateVal)

			columns.SetTitle(getCoords.Results[0].FormattedAddress)

			for i := 0; i <= 12; i++ {
				columns.SetCell(i, 0, tview.NewTableCell(getForecast.Properties.Periods[i].Name).SetTextColor(tcell.ColorBlue).SetAlign(tview.AlignCenter))
				
				columns.SetCell(i, 1, tview.NewTableCell(strconv.Itoa(getForecast.Properties.Periods[i].Temperature) + "°F").SetTextColor(tcell.ColorRed).SetAlign(tview.AlignCenter))
				
				columns.SetCell(i, 2, tview.NewTableCell(getForecast.Properties.Periods[i].WindSpeed).SetTextColor(tcell.ColorRed).SetAlign(tview.AlignCenter))
				
				columns.SetCell(i, 3, tview.NewTableCell(getForecast.Properties.Periods[i].ShortForecast).SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
				
			}

			weatherAppUI.SetFocus(columns)
		})

	flex = tview.NewFlex().
		AddItem(citySelection, 0, 1, true).
		AddItem(columns, 0, 3, false)

	weatherAppUI.SetRoot(flex, true).SetFocus(citySelection)
	weatherAppUI.Run()

}

func apiCalls(city string, state string) (GeocodeResponse, ForecastResponse) {
	googleResponse := googleQuery(city, state)

	latCnv := strconv.FormatFloat(googleResponse.Results[0].Geometry.Location.Lat, 'f', 5, 64)
	lngCnv := strconv.FormatFloat(googleResponse.Results[0].Geometry.Location.Lng, 'f', 5, 64)

	pointResponse := weatherPointQuery(latCnv, lngCnv)

	weatherForecastResponse := forecastQuery(pointResponse.Properties.Forecast)

	return googleResponse, weatherForecastResponse
}

func googleQuery(city string, state string) GeocodeResponse {
	googleApiEndpoint := "https://maps.googleapis.com/maps/api/geocode/json?address="
	googleApiKey := "&key=CHANGEME"

	reqEndpoint := googleApiEndpoint + city + "+" + state + googleApiKey

	googleResponse := makeRequest(reqEndpoint)

	var responseObj GeocodeResponse
	json.Unmarshal(googleResponse, &responseObj)

	return responseObj
}

func weatherPointQuery(lat string, lng string) WeatherPointResponse {
	weatherPtApiEndpoint := "https://api.weather.gov/points/"
	reqEndpoint := weatherPtApiEndpoint + lat + "," + lng

	pointResponse := makeRequest(reqEndpoint)

	var responseObj WeatherPointResponse
	json.Unmarshal(pointResponse, &responseObj)

	return responseObj
}

func forecastQuery(url string) ForecastResponse {
	forecastApiEndpoint := url

	forecastResponse := makeRequest(forecastApiEndpoint)

	var responseObj ForecastResponse
	json.Unmarshal(forecastResponse, &responseObj)

	return responseObj
}

func makeRequest(endpoint string) []byte {
	client := &http.Client{}
	request, err := http.NewRequest("GET", endpoint, nil)
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

	return responseBytes
}
