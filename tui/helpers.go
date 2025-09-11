package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/esferadigital/clima/openmeteo"
)

func GetDefaultForecast() tea.Msg {
	params := openmeteo.ForecastParams{
		Latitude:  DEFAULT_LATITUDE,
		Longitude: DEFAULT_LONGITUDE,
		Current: []openmeteo.CurrentWeatherVariables{
			openmeteo.Temperature2m,
			openmeteo.RelativeHumidity2m,
			openmeteo.ApparentTemperature,
			openmeteo.IsDay,
			openmeteo.WeatherCode,
			openmeteo.WindSpeed10m,
			openmeteo.WindDirection10m,
			openmeteo.WindGusts10m,
			openmeteo.Rain,
			openmeteo.Showers,
			openmeteo.Snowfall,
			openmeteo.CloudCover,
			openmeteo.SeaLevelPressure,
			openmeteo.SurfacePressure,
			openmeteo.PressureAtGround,
		},
	}

	forecast, err := openmeteo.GetForecast(params)
	if err != nil {
		return err
	}

	return forecast
}

