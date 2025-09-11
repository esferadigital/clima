package openmeteo

import "strings"

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
func WriteVariableCSV(variables []CurrentWeatherVariables) string {
	variableNames := make([]string, len(variables))
	for i, variable := range variables {
		variableNames[i] = string(variable)
	}
	return strings.Join(variableNames, ",")
}

