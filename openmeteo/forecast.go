package openmeteo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Parameters for the Open-Meteo Forecast V1 API.
// These are not exclusive. Check the docs for additional ones.
// https://open-meteo.com/en/docs
type ForecastParams struct {
	Latitude  float64
	Longitude float64
	Current   []CurrentWeatherVariables
}

// Response from the Open-Meteo Forecast V1 API.
// For possible values of `Current`,
// refer to the `openmeteo/variables.go` file
// for possible current weather variables and their units.
type ForecastResponse struct {
	Latitude         float64        `json:"latitude"`
	Longitude        float64        `json:"longitude"`
	Elevation        float64        `json:"elevation"`
	GenerationTimeMs float64        `json:"generation_time_ms"`
	UTCOffsetSeconds int            `json:"utc_offset_seconds"`
	Timezone         string         `json:"timezone"`
	TimezoneAbbrev   string         `json:"timezone_abbreviation"`
	CurrentUnits     map[string]any `json:"current_units"`
	Current          map[string]any `json:"current"`
}

// Variables available to request from the Open-Meteo Forecast V1 API.
type CurrentWeatherVariables string

const (
	Temperature2m       CurrentWeatherVariables = "temperature_2m"
	RelativeHumidity2m  CurrentWeatherVariables = "relative_humidity_2m"
	ApparentTemperature CurrentWeatherVariables = "apparent_temperature"
	IsDay               CurrentWeatherVariables = "is_day"
	WeatherCode         CurrentWeatherVariables = "weather_code"
	CloudCover          CurrentWeatherVariables = "cloud_cover"
	SeaLevelPressure    CurrentWeatherVariables = "pressure_msl"
	SurfacePressure     CurrentWeatherVariables = "surface_pressure"
	PressureAtGround    CurrentWeatherVariables = "precipitation"
	Rain                CurrentWeatherVariables = "rain"
	Showers             CurrentWeatherVariables = "showers"
	Snowfall            CurrentWeatherVariables = "snowfall"
	WindSpeed10m        CurrentWeatherVariables = "wind_speed_10m"
	WindDirection10m    CurrentWeatherVariables = "wind_direction_10m"
	WindGusts10m        CurrentWeatherVariables = "wind_gusts_10m"
)

// Compose a comma-separated string of variable names.
// This is used to send the variables as a query parameter.
func writeVariableCSV(variables []CurrentWeatherVariables) string {
	variableNames := make([]string, len(variables))
	for i, variable := range variables {
		variableNames[i] = string(variable)
	}
	return strings.Join(variableNames, ",")
}

const FORECAST_API_URL = "https://api.open-meteo.com/v1/forecast"

// Retrieve the current forecast data for a given location and parameters.
// Data is provided by the Open-Meteo API.
func GetForecast(params ForecastParams) (ForecastResponse, error) {
	variables := writeVariableCSV(params.Current)
	url := fmt.Sprintf("%s?latitude=%f&longitude=%f&current=%s", FORECAST_API_URL, params.Latitude, params.Longitude, variables)
	resp, err := http.Get(url)
	if err != nil {
		return ForecastResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ForecastResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	var response ForecastResponse
	err = decoder.Decode(&response)
	if err != nil {
		return ForecastResponse{}, err
	}

	return response, nil
}

