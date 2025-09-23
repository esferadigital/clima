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
	Daily     []DailyWeatherVariables
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
	DailyUnits       map[string]any `json:"daily_units"`
	Daily            map[string]any `json:"daily"`
}

// Variables available to request from the Open-Meteo Forecast V1 API for current weather.
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
	Precipitation       CurrentWeatherVariables = "precipitation"
	Rain                CurrentWeatherVariables = "rain"
	Showers             CurrentWeatherVariables = "showers"
	Snowfall            CurrentWeatherVariables = "snowfall"
	WindSpeed10m        CurrentWeatherVariables = "wind_speed_10m"
	WindDirection10m    CurrentWeatherVariables = "wind_direction_10m"
	WindGusts10m        CurrentWeatherVariables = "wind_gusts_10m"
)

// Variables available to request from the Open-Meteo Forecast V1 API for daily weather.
type DailyWeatherVariables string

const (
	Temperature2mMin DailyWeatherVariables = "temperature_2m_min"
	Temperature2mMax DailyWeatherVariables = "temperature_2m_max"
	UVIndexMax       DailyWeatherVariables = "uv_index_max"
)

// Compose a comma-separated string of variable names.
// This is used to send the variables as a query parameter.
func writeVariableCSV[T ~string](variables []T) string {
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
	url := fmt.Sprintf("%s?latitude=%f&longitude=%f", FORECAST_API_URL, params.Latitude, params.Longitude)
	if len(params.Current) > 0 {
		currentVars := writeVariableCSV(params.Current)
		url += fmt.Sprintf("&current=%s", currentVars)
	}
	if len(params.Daily) > 0 {
		dailyVars := writeVariableCSV(params.Daily)
		url += fmt.Sprintf("&daily=%s", dailyVars)
	}

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

func MapWeatherCode(code float64) string {
	wmoCodes := map[float64]string{
		0:  "Clear",
		1:  "Mostly clear",
		2:  "Partly cloudy",
		3:  "Overcast",
		45: "Fog",
		48: "Icy fog",
		51: "Light drizzle",
		53: "Drizzle",
		55: "Heavy drizzle",
		56: "Light freezing drizzle",
		57: "Heavy freezing drizzle",
		61: "Light rain",
		63: "Rain",
		65: "Heavy rain",
		66: "Light freezing rain",
		67: "Heavy freezing rain",
		71: "Light snowfall",
		73: "Snowfall",
		75: "Heavy snowfall",
		77: "Snow grains",
		80: "Light rain showers",
		81: "Rain showers",
		82: "Heavy rain showers",
		85: "Light snow showers",
		86: "Heavy snow showers",
		95: "Thunderstorm",
		96: "Thunderstorm with light hail",
		99: "Thunderstorm with heavy hail",
	}

	return wmoCodes[code]
}
