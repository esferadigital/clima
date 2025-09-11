package openmeteo

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

