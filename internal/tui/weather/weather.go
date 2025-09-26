package weather

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/esferadigital/clima/internal/openmeteo"
	"github.com/esferadigital/clima/internal/store"
)

// ---- styles ----

var (
	accent = lipgloss.NewStyle().Foreground(lipgloss.Color("13"))
	subtle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	label  = lipgloss.NewStyle().Width(20).Foreground(lipgloss.Color("8"))
)

// ---- keymap ----

type keyMap struct {
	newSearch       key.Binding
	recentLocations key.Binding
	refresh         key.Binding
	quit            key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.newSearch, k.recentLocations, k.refresh, k.quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.newSearch}, {k.recentLocations},
		{k.refresh}, {k.quit},
	}
}

func newKeyMap() keyMap {
	return keyMap{
		newSearch: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new search"),
		),
		recentLocations: key.NewBinding(
			key.WithKeys("b"),
			key.WithHelp("b", "recent locations"),
		),
		refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
	}
}

// ---- msg ----

type dataMsg struct {
	forecast openmeteo.ForecastResponse
}

type errorMsg struct {
	err error
}

type savedMsg struct {
	err error
}

type NewSearchMsg struct{}

type RecentMsg struct{}

// ---- cmd ----

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
			return errorMsg{
				err: err,
			}
		}
		return dataMsg{
			forecast: res,
		}
	}
}

func requestNewSearchCmd() tea.Cmd {
	return func() tea.Msg {
		return NewSearchMsg{}
	}
}

func requestRecentCmd() tea.Cmd {
	return func() tea.Msg {
		return RecentMsg{}
	}
}

func saveRecentLocationCmd(location openmeteo.GeocodingResult) tea.Cmd {
	return func() tea.Msg {
		err := store.AddRecentLocation(location)
		return savedMsg{err: err}
	}
}

// ---- model ----

type view int

const (
	viewLoading = iota
	viewReady
	viewError
)

type Model struct {
	view     view
	errStr   string
	ellipsis spinner.Model
	location openmeteo.GeocodingResult
	forecast openmeteo.ForecastResponse
	keys     keyMap
	help     help.Model
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		saveRecentLocationCmd(m.location),
		getForecastCmd(m.location.Latitude, m.location.Longitude),
		m.ellipsis.Tick,
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keys.newSearch) {
			return m, requestNewSearchCmd()
		}
		if key.Matches(msg, m.keys.recentLocations) {
			return m, requestRecentCmd()
		}
		if key.Matches(msg, m.keys.refresh) && m.view == viewReady {
			m.view = viewLoading
			return m, tea.Batch(getForecastCmd(m.location.Latitude, m.location.Longitude), m.ellipsis.Tick)
		}
		if key.Matches(msg, m.keys.quit) {
			return m, tea.Quit
		}
	case dataMsg:
		m.forecast = msg.forecast
		m.view = viewReady
		return m, nil
	case errorMsg:
		m.view = viewError
		m.errStr = msg.err.Error()
		return m, nil
	}

	if m.view == viewLoading {
		var cmd tea.Cmd
		m.ellipsis, cmd = m.ellipsis.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	switch m.view {
	case viewLoading:
		return fmt.Sprintf("\nLoading forecast%s\n", m.ellipsis.View())
	case viewReady:
		weather := m.forecast
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
				minTempValue = "-"
			}
		} else {
			minTempValue = "-"
		}
		s += minTempLabel + minTempValue

		maxTempLabel := label.Render("\nMax")
		var maxTempValue string
		if maxArray, ok := weather.Daily[string(openmeteo.Temperature2mMax)].([]any); ok && len(maxArray) > 0 {
			if maxTemp, ok := maxArray[0].(float64); ok {
				maxTempValue = fmt.Sprintf("%.1f %s", maxTemp, weather.DailyUnits[string(openmeteo.Temperature2mMax)])
			} else {
				maxTempValue = "-"
			}
		} else {
			maxTempValue = "-"
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

		helpView := m.help.View(m.keys)
		return s + "\n\n" + helpView
	case viewError:
		return "\nFailed to get weather forecast:" + m.errStr
	default:
		return "\nunknown error state (weather)"
	}
}

func New(location openmeteo.GeocodingResult) Model {
	ellipsis := spinner.New()
	ellipsis.Spinner = spinner.Ellipsis
	ellipsis.Style = accent

	return Model{
		view:     viewLoading,
		location: location,
		ellipsis: ellipsis,
		keys:     newKeyMap(),
		help:     help.New(),
	}
}
