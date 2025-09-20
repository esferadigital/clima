package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/esferadigital/clima/openmeteo"
)

// LocationItem wraps GeocodingResult to implement list.Item interface
type LocationItem struct {
	openmeteo.GeocodingResult
}

// FilterValue implements the list.Item interface
func (l LocationItem) FilterValue() string {
	return l.Name
}

// Title returns the display title for the list item
func (l LocationItem) Title() string {
	return fmt.Sprintf("%s, %s", l.Name, l.Country)
}

// Description returns the display description for the list item
func (l LocationItem) Description() string {
	return fmt.Sprintf("Lat: %.4f, Lon: %.4f", l.Latitude, l.Longitude)
}

// ---- status; part of model, and used to conditionally render content ----

type ModelStatus int

const (
	// search location
	locationSearch = iota
	locationLoading
	locationPick

	// forecast
	forecastLoading
	forecastReady

	failed
)

// ---- messages ----

type locationsFoundMsg struct {
	locations []openmeteo.GeocodingResult
}

type locationSelectedMsg struct {
	location openmeteo.GeocodingResult
}

type locationErrorMsg struct {
	err error
}

type forecastLoadedMsg struct {
	forecast openmeteo.ForecastResponse
}

type forecastErrorMsg struct {
	err error
}

// ---- model ----

type Model struct {
	sink          io.Writer
	status        ModelStatus
	locationInput textinput.Model
	locations     []openmeteo.GeocodingResult
	locationList  list.Model
	location      openmeteo.GeocodingResult
	weather       openmeteo.ForecastResponse
	err           string
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.sink != nil {
		fmt.Fprintf(m.sink, "[MSG] %T: %+v\n", msg, msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			if m.status == locationSearch {
				m.status = locationLoading
				return m, getLocationsCmd(m.locationInput.Value())
			}

			if m.status == locationPick {
				picked, ok := m.locationList.SelectedItem().(LocationItem)
				if ok {
					return m, selectLocationCmd(picked.GeocodingResult)
				}
			}
		}

		if msg.Type == tea.KeyCtrlC || (msg.String() == "q" && m.status != locationSearch) {
			return m, tea.Quit
		}

	case locationsFoundMsg:
		m.locations = msg.locations

		// Auto-select single result
		if len(msg.locations) == 1 {
			m.location = msg.locations[0]
			m.status = forecastLoading
			return m, getForecastCmd(m.location.Latitude, m.location.Longitude)
		}

		// Show list for multiple results
		items := make([]list.Item, len(msg.locations))
		for i, loc := range msg.locations {
			items[i] = LocationItem{loc}
		}
		m.locationList.SetItems(items)
		m.status = locationPick
		return m, nil

	case locationSelectedMsg:
		m.location = msg.location
		m.status = forecastLoading
		return m, getForecastCmd(m.location.Latitude, m.location.Longitude)

	case locationErrorMsg:
		m.err = "Failed to find location: " + msg.err.Error()
		m.status = failed
		return m, nil

	case forecastLoadedMsg:
		m.weather = msg.forecast
		m.status = forecastReady
		return m, tea.Quit

	case forecastErrorMsg:
		m.err = "Failed to get forecast: " + msg.err.Error()
		m.status = failed
		return m, nil
	}

	var cmd tea.Cmd
	switch m.status {
	case locationSearch:
		m.locationInput, cmd = m.locationInput.Update(msg)
	case locationPick:
		m.locationList, cmd = m.locationList.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	switch m.status {
	case failed:
		return "\n" + m.err + "\n\nPress 'q' to quit."
	case locationSearch:
		return "\nLocation search:\n" + m.locationInput.View()
	case locationLoading:
		return "\nLoading location ...\n"
	case locationPick:
		return "\nPick a location:\n\n" + m.locationList.View()
	case forecastLoading:
		return "\nforecast loading ...\n"
	case forecastReady:
		weather := m.weather
		s := ""

		temperature := fmt.Sprintf("\nTemperature: %.1f %s", weather.Current[string(openmeteo.Temperature2m)], weather.CurrentUnits[string(openmeteo.Temperature2m)])
		s += temperature

		windSpeed := fmt.Sprintf("\nWind Speed: %.1f %s", weather.Current[string(openmeteo.WindSpeed10m)], weather.CurrentUnits[string(openmeteo.WindSpeed10m)])
		s += windSpeed

		windDirection := fmt.Sprintf("\nWind Direction: %.1f %s", weather.Current[string(openmeteo.WindDirection10m)], weather.CurrentUnits[string(openmeteo.WindDirection10m)])
		s += windDirection

		windGusts := fmt.Sprintf("\nWind Gusts: %.1f %s", weather.Current[string(openmeteo.WindGusts10m)], weather.CurrentUnits[string(openmeteo.WindGusts10m)])
		s += windGusts

		rain := fmt.Sprintf("\nRain: %.1f %s", weather.Current[string(openmeteo.Rain)], weather.CurrentUnits[string(openmeteo.Rain)])
		s += rain

		showers := fmt.Sprintf("\nShowers: %.1f %s", weather.Current[string(openmeteo.Showers)], weather.CurrentUnits[string(openmeteo.Showers)])
		s += showers

		snowfall := fmt.Sprintf("\nSnowfall: %.1f %s", weather.Current[string(openmeteo.Snowfall)], weather.CurrentUnits[string(openmeteo.Snowfall)])
		s += snowfall

		cloudCover := fmt.Sprintf("\nCloud Cover: %.1f %s", weather.Current[string(openmeteo.CloudCover)], weather.CurrentUnits[string(openmeteo.CloudCover)])
		s += cloudCover

		seaLevelPressure := fmt.Sprintf("\nSea Level Pressure: %.1f %s", weather.Current[string(openmeteo.SeaLevelPressure)], weather.CurrentUnits[string(openmeteo.SeaLevelPressure)])
		s += seaLevelPressure

		surfacePressure := fmt.Sprintf("\nSurface Pressure: %.1f %s", weather.Current[string(openmeteo.SurfacePressure)], weather.CurrentUnits[string(openmeteo.SurfacePressure)])
		s += surfacePressure

		pressureAtGround := fmt.Sprintf("\nPressure at Ground: %.1f %s", weather.Current[string(openmeteo.PressureAtGround)], weather.CurrentUnits[string(openmeteo.PressureAtGround)])
		s += pressureAtGround

		return s

	default:
		return ""
	}
}

func InitialModel(sink io.Writer) Model {
	locationInput := textinput.New()
	locationInput.Placeholder = "Salinas"
	locationInput.Focus()
	locationInput.CharLimit = 256
	locationInput.Width = 20

	listDelegate := list.NewDefaultDelegate()
	listDelegate.ShowDescription = false

	locationList := list.New([]list.Item{}, listDelegate, 30, 14)
	locationList.SetShowStatusBar(false)
	locationList.SetFilteringEnabled(false)
	locationList.SetShowHelp(false)
	locationList.SetShowTitle(false)

	return Model{
		sink:          sink,
		locationInput: locationInput,
		locationList:  locationList,
		status:        locationSearch,
	}
}
