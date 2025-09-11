package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/esferadigital/clima/openmeteo"
)

type ModelStatus int

const (
	statusLoading ModelStatus = iota
	statusLoaded
	statusFailed
)

type Model struct {
	status ModelStatus
	weather openmeteo.ForecastResponse
	err error
}

func InitialModel() Model {
	return Model{ status: statusLoading }
}

func (m Model) Init() tea.Cmd {
	return GetDefaultForecast
}


func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case openmeteo.ForecastResponse:
		m.weather = openmeteo.ForecastResponse(msg)
		m.status = statusLoaded
		return m, tea.Quit

	case error:
		m.status = statusFailed
		m.err = msg
		return m, tea.Quit

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	switch m.status {
	case statusLoading:
		return "\nGetting the weather ... (press q to quit)\n"
	case statusFailed:
		return fmt.Sprintf("\nFailed to get the weather: %v\n", m.err)
	case statusLoaded:
		forecast := m.weather
		s := ""

		temperature := fmt.Sprintf("\nTemperature: %.1f %s", forecast.Current[string(openmeteo.Temperature2m)], forecast.CurrentUnits[string(openmeteo.Temperature2m)])
		s += temperature

		windSpeed := fmt.Sprintf("\nWind Speed: %.1f %s", forecast.Current[string(openmeteo.WindSpeed10m)], forecast.CurrentUnits[string(openmeteo.WindSpeed10m)])
		s += windSpeed

		windDirection := fmt.Sprintf("\nWind Direction: %.1f %s", forecast.Current[string(openmeteo.WindDirection10m)], forecast.CurrentUnits[string(openmeteo.WindDirection10m)])
		s += windDirection

		windGusts := fmt.Sprintf("\nWind Gusts: %.1f %s", forecast.Current[string(openmeteo.WindGusts10m)], forecast.CurrentUnits[string(openmeteo.WindGusts10m)])
		s += windGusts

		rain := fmt.Sprintf("\nRain: %.1f %s", forecast.Current[string(openmeteo.Rain)], forecast.CurrentUnits[string(openmeteo.Rain)])
		s += rain

		showers := fmt.Sprintf("\nShowers: %.1f %s", forecast.Current[string(openmeteo.Showers)], forecast.CurrentUnits[string(openmeteo.Showers)])
		s += showers

		snowfall := fmt.Sprintf("\nSnowfall: %.1f %s", forecast.Current[string(openmeteo.Snowfall)], forecast.CurrentUnits[string(openmeteo.Snowfall)])
		s += snowfall

		cloudCover := fmt.Sprintf("\nCloud Cover: %.1f %s", forecast.Current[string(openmeteo.CloudCover)], forecast.CurrentUnits[string(openmeteo.CloudCover)])
		s += cloudCover

		seaLevelPressure := fmt.Sprintf("\nSea Level Pressure: %.1f %s", forecast.Current[string(openmeteo.SeaLevelPressure)], forecast.CurrentUnits[string(openmeteo.SeaLevelPressure)])
		s += seaLevelPressure

		surfacePressure := fmt.Sprintf("\nSurface Pressure: %.1f %s", forecast.Current[string(openmeteo.SurfacePressure)], forecast.CurrentUnits[string(openmeteo.SurfacePressure)])
		s += surfacePressure

		pressureAtGround := fmt.Sprintf("\nPressure at Ground: %.1f %s", forecast.Current[string(openmeteo.PressureAtGround)], forecast.CurrentUnits[string(openmeteo.PressureAtGround)])
		s += pressureAtGround

		return s
	
	default:
		return ""
	}
}

