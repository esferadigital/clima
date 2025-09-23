package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/esferadigital/clima/openmeteo"
)

// ---- commands ----

const DEFAULT_SEARCH_COUNT = 10

func getLocationsCmd() tea.Cmd {
	return func() tea.Msg {
		locations, err := loadRecentLocations()
		if err != nil {
			return recentLocationsLoadedMsg{
				locations: []openmeteo.GeocodingResult{},
			}
		}

		return recentLocationsLoadedMsg{
			locations: locations,
		}
	}
}

func saveRecentLocationCmd(location openmeteo.GeocodingResult) tea.Cmd {
	return func() tea.Msg {
		err := addRecentLocation(location)
		return recentLocationSavedMsg{err: err}
	}
}

func searchLocationsCmd(name string) tea.Cmd {
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
				openmeteo.ApparentTemperature,
				openmeteo.RelativeHumidity2m,
				openmeteo.IsDay,
				openmeteo.WeatherCode,
				openmeteo.WindSpeed10m,
				openmeteo.WindDirection10m,
				openmeteo.WindGusts10m,
				openmeteo.Precipitation,
				openmeteo.SeaLevelPressure,
			},
			Daily: []openmeteo.DailyWeatherVariables{
				openmeteo.Temperature2mMin,
				openmeteo.Temperature2mMax,
				openmeteo.UVIndexMax,
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
