package openmeteo

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const FORECAST_API_URL = "https://api.open-meteo.com/v1/forecast"

// Retrieve the current forecast data for a given location and parameters.
// Data is provided by the Open-Meteo API.
func GetForecast(params ForecastParams) (*ForecastResponse, error) {
	variables := WriteVariableCSV(params.Current)
	url := fmt.Sprintf("%s?latitude=%f&longitude=%f&current=%s", FORECAST_API_URL, params.Latitude, params.Longitude, variables)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	var response ForecastResponse
	err = decoder.Decode(&response)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return nil, err
	}

	return &response, nil
}
