package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/esferadigital/clima/openmeteo"
)

var (
	accent = lipgloss.NewStyle().Foreground(lipgloss.Color("13"))
	subtle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	label = lipgloss.NewStyle().Width(20).Foreground(lipgloss.Color("8"))
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
	sink            io.Writer
	status          ModelStatus
	locationSpinner spinner.Model
	forecastSpinner spinner.Model
	locationInput   textinput.Model
	locations       []openmeteo.GeocodingResult
	locationList    list.Model
	location        openmeteo.GeocodingResult
	weather         openmeteo.ForecastResponse
	err             string
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
				return m, tea.Batch(getLocationsCmd(m.locationInput.Value()), m.locationSpinner.Tick)
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
			return m, tea.Batch(getForecastCmd(m.location.Latitude, m.location.Longitude), m.forecastSpinner.Tick)
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
		return m, tea.Batch(getForecastCmd(m.location.Latitude, m.location.Longitude), m.forecastSpinner.Tick)

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
	case locationLoading:
		m.locationSpinner, cmd = m.locationSpinner.Update(msg)
	case locationPick:
		m.locationList, cmd = m.locationList.Update(msg)
	case forecastLoading:
		m.forecastSpinner, cmd = m.forecastSpinner.Update(msg)
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
		return fmt.Sprintf("\nFinding location%s\n", m.locationSpinner.View())
	case locationPick:
		return "\nPick a location:\n\n" + m.locationList.View()
	case forecastLoading:
		return fmt.Sprintf("\nLoading forecast%s\n", m.forecastSpinner.View())
	case forecastReady:
		weather := m.weather
		s := fmt.Sprintf("\n%s, %s", m.location.Name, m.location.Country)

		weatherCode, ok := weather.Current[string(openmeteo.WeatherCode)].(float64)
		if ok {
			weatherInterpretation := fmt.Sprintf("\n%s", openmeteo.MapWeatherCode(weatherCode))
			s += accent.Render(weatherInterpretation)
		}

		temperature := fmt.Sprintf("\n%.1f %s", weather.Current[string(openmeteo.Temperature2m)], weather.CurrentUnits[string(openmeteo.Temperature2m)])
		temperature = temperature + subtle.Render(fmt.Sprintf(" (feels like %.1f %s)", weather.Current[string(openmeteo.ApparentTemperature)], weather.CurrentUnits[string(openmeteo.ApparentTemperature)]))
		s += temperature

		minTempLabel := label.Render("\nMin")
		var minTempValue string
		if minArray, ok := weather.Daily[string(openmeteo.Temperature2mMin)].([]any); ok && len(minArray) > 0 {
			if minTemp, ok := minArray[0].(float64); ok {
				minTempValue = fmt.Sprintf("%.1f %s", minTemp, weather.DailyUnits[string(openmeteo.Temperature2mMin)])
			} else {
				minTempValue = "N/A"
			}
		} else {
			minTempValue = "N/A"
		}
		s += minTempLabel + minTempValue

		maxTempLabel := label.Render("\nMax")
		var maxTempValue string
		if maxArray, ok := weather.Daily[string(openmeteo.Temperature2mMax)].([]any); ok && len(maxArray) > 0 {
			if maxTemp, ok := maxArray[0].(float64); ok {
				maxTempValue = fmt.Sprintf("%.1f %s", maxTemp, weather.DailyUnits[string(openmeteo.Temperature2mMax)])
			} else {
				maxTempValue = "N/A"
			}
		} else {
			maxTempValue = "N/A"
		}
		s += maxTempLabel + maxTempValue

		windLabel := label.Render("\n\nWind")
		windValue := fmt.Sprintf("%.1f %s @ %.1f %s", weather.Current[string(openmeteo.WindSpeed10m)], weather.CurrentUnits[string(openmeteo.WindSpeed10m)], weather.Current[string(openmeteo.WindDirection10m)], weather.CurrentUnits[string(openmeteo.WindDirection10m)])
		s += windLabel + windValue

		windGustsLabel := label.Render("\nWind gusts")
		windGustsValue := fmt.Sprintf("%.1f %s", weather.Current[string(openmeteo.WindGusts10m)], weather.CurrentUnits[string(openmeteo.WindGusts10m)])
		s += windGustsLabel + windGustsValue

		humidityLabel := label.Render("\nHumidity")
		humidityValue := fmt.Sprintf("%.1f %s", weather.Current[string(openmeteo.RelativeHumidity2m)], weather.CurrentUnits[string(openmeteo.RelativeHumidity2m)])
		s += humidityLabel + humidityValue

		precipitationLabel := label.Render("\nPrecipitation")
		precipitationValue := fmt.Sprintf("%.1f %s", weather.Current[string(openmeteo.Precipitation)], weather.CurrentUnits[string(openmeteo.Precipitation)])
		s += precipitationLabel + precipitationValue

		pressureLabel := label.Render("\nPressure")
		pressureValue := fmt.Sprintf("%.1f %s", weather.Current[string(openmeteo.SeaLevelPressure)], weather.CurrentUnits[string(openmeteo.SeaLevelPressure)])
		s += pressureLabel + pressureValue

		uvLabel := label.Render("\nUV index")
		var uvValue string
		if uvArray, ok := weather.Daily[string(openmeteo.UVIndexMax)].([]any); ok && len(uvArray) > 0 {
			if uvToday, ok := uvArray[0].(float64); ok {
				uvValue = fmt.Sprintf("%.1f", uvToday)
			} else {
				uvValue = "-"
			}
		} else {
			uvValue = "-"
		}
		s += uvLabel + uvValue

		return s + "\n"

	default:
		return ""
	}
}

func InitialModel(sink io.Writer) Model {
	locationSpin := spinner.New()
	locationSpin.Spinner = spinner.Ellipsis
	locationSpin.Style = accent

	forecastSpin := spinner.New()
	forecastSpin.Spinner = spinner.Ellipsis
	forecastSpin.Style = accent

	locationInput := textinput.New()
	locationInput.Placeholder = "Salinas"
	locationInput.Focus()
	locationInput.CharLimit = 256
	locationInput.Width = 20
	locationInput.Cursor.Style = accent

	listDelegate := list.NewDefaultDelegate()
	listDelegate.ShowDescription = false

	locationList := list.New([]list.Item{}, listDelegate, 30, 14)
	locationList.SetShowStatusBar(false)
	locationList.SetFilteringEnabled(false)
	locationList.SetShowHelp(true)
	locationList.SetShowTitle(false)

	return Model{
		sink:            sink,
		locationSpinner: locationSpin,
		forecastSpinner: forecastSpin,
		locationInput:   locationInput,
		locationList:    locationList,
		status:          locationSearch,
	}
}
