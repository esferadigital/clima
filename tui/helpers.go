package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/esferadigital/clima/openmeteo"
)

const DEFAULT_SEARCH_COUNT = 10

func getLocationsCmd(name string) tea.Cmd {
	return func() tea.Msg {
		params := openmeteo.GeocodingParams{
			Name:  name,
			Count: DEFAULT_SEARCH_COUNT,
		}

		res, err := openmeteo.SearchLocation(params)
		if err != nil {
			return locationErrorMsg{
				err: err,
			}
		}

		return locationsFoundMsg{
			locations: res.Results,
		}
	}
}

func selectLocationCmd(selected openmeteo.GeocodingResult) tea.Cmd {
	return func() tea.Msg {
		return locationSelectedMsg{location: selected}
	}
}

func getForecastCmd(lat float64, long float64) tea.Cmd {
	return func() tea.Msg {
		params := openmeteo.ForecastParams{
			Latitude:  lat,
			Longitude: long,
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

		res, err := openmeteo.GetForecast(params)
		if err != nil {
			return forecastErrorMsg{
				err: err,
			}
		}

		return forecastLoadedMsg{
			forecast: res,
		}
	}
}

