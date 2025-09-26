package search

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/esferadigital/clima/internal/openmeteo"
)

const DEFAULT_SEARCH_COUNT = 10

// ---- styles ----

var accent = lipgloss.NewStyle().Foreground(lipgloss.Color("13"))

// ---- keymap ----

type inputKeyMap struct {
	submit     key.Binding
	exitSearch key.Binding
}

func (k inputKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.submit, k.exitSearch}
}

func (k inputKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.submit, k.exitSearch},
	}
}

func newInputKeyMap() inputKeyMap {
	return inputKeyMap{
		submit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "search"),
		),
		exitSearch: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "exit search"),
		),
	}
}

type listKeyMap struct {
	up        key.Binding
	down      key.Binding
	pick      key.Binding
	newSearch key.Binding
	quit      key.Binding
}

func (k listKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.up, k.down, k.pick, k.newSearch, k.quit}
}

func (k listKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.up},
		{k.down},
		{k.pick},
		{k.newSearch},
		{k.quit},
	}
}

func newListKeyMap() listKeyMap {
	return listKeyMap{
		up: key.NewBinding(
			key.WithKeys("↑"),
			key.WithHelp("↑", "up"),
		),
		down: key.NewBinding(
			key.WithKeys("↓"),
			key.WithHelp("↓", "down"),
		),
		pick: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "pick"),
		),
		newSearch: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new search"),
		),
		quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
	}
}

// ---- msg ----

type dataMsg struct {
	locations []openmeteo.GeocodingResult
}

type errorMsg struct {
	err error
}

type SearchCompleteMsg struct {
	Location openmeteo.GeocodingResult
}

type RecentMsg struct{}

// ---- helpers ----

// Implements list.Item interface and wraps openmeteo.GeocodingResult
type searchListItem struct {
	openmeteo.GeocodingResult
}

func (i searchListItem) FilterValue() string {
	return i.Name
}

func (i searchListItem) Title() string {
	return fmt.Sprintf("%s, %s", i.Name, i.Country)
}

func (i searchListItem) Description() string {
	return fmt.Sprintf("Lat: %.4f, Lon: %.4f", i.Latitude, i.Longitude)
}

// ---- cmd ----

func searchLocationsCmd(name string) tea.Cmd {
	return func() tea.Msg {
		params := openmeteo.GeocodingParams{
			Name:  name,
			Count: DEFAULT_SEARCH_COUNT,
		}
		res, err := openmeteo.SearchLocation(params)
		if err != nil {
			return errorMsg{
				err: err,
			}
		}
		return dataMsg{
			locations: res.Results,
		}
	}
}

func pickCmd(location openmeteo.GeocodingResult) tea.Cmd {
	return func() tea.Msg {
		return SearchCompleteMsg{
			Location: location,
		}
	}
}

func requestRecentCmd() tea.Cmd {
	return func() tea.Msg {
		return RecentMsg{}
	}
}

// ---- model ----

type view int

const (
	viewInput = iota
	viewLoading
	viewPick
	viewError
)

type Model struct {
	view      view
	input     textinput.Model
	inputKeys inputKeyMap
	ellipsis  spinner.Model
	list      list.Model
	listKeys  listKeyMap
	help      help.Model
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.view == viewInput {
			if key.Matches(msg, m.inputKeys.submit) {
				m.view = viewLoading
				return m, tea.Batch(searchLocationsCmd(m.input.Value()), m.ellipsis.Tick)
			}
			if key.Matches(msg, m.inputKeys.exitSearch) {
				return m, requestRecentCmd()
			}
		}
		if m.view == viewPick {
			if key.Matches(msg, m.inputKeys.submit) {
				picked, ok := m.list.SelectedItem().(searchListItem)
				if ok {
					return m, pickCmd(picked.GeocodingResult)
				}
			}
			if key.Matches(msg, m.listKeys.newSearch) {
				m.input.Reset()
				m.view = viewInput
				return m, nil
			}
		}

	case dataMsg:
		if len(msg.locations) == 1 {
			return m, pickCmd(msg.locations[0])
		}
		items := make([]list.Item, len(msg.locations))
		for i, loc := range msg.locations {
			items[i] = searchListItem{loc}
		}
		m.list.SetItems(items)
		m.view = viewPick
		return m, nil
	case errorMsg:
		m.view = viewError
		return m, nil
	}

	// Forward messages to sub-components
	var cmd tea.Cmd
	switch m.view {
	case viewInput:
		m.input, cmd = m.input.Update(msg)
	case viewLoading:
		m.ellipsis, cmd = m.ellipsis.Update(msg)
	case viewPick:
		m.list, cmd = m.list.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	view := "\n"
	switch m.view {
	case viewInput:
		return view + "Location search:\n" + m.input.View() + "\n\n" + m.help.View(m.inputKeys)
	case viewLoading:
		return view + fmt.Sprintf("Finding location%s\n", m.ellipsis.View())
	case viewPick:
		return view + "Pick a location:\n\n" + m.list.View() + "\n" + m.help.View(m.listKeys)
	case viewError:
		return view + "Error ... try again or quit\n"
	default:
		return ""
	}
}

func New() Model {
	input := textinput.New()
	input.Placeholder = "Salinas"
	input.Focus()
	input.CharLimit = 256
	input.Width = 20
	input.Cursor.Style = accent

	ellipsis := spinner.New()
	ellipsis.Spinner = spinner.Ellipsis
	ellipsis.Style = accent

	listDelegate := list.NewDefaultDelegate()
	listDelegate.ShowDescription = false
	list := list.New([]list.Item{}, listDelegate, 30, 14)
	list.SetShowStatusBar(false)
	list.SetFilteringEnabled(false)
	list.SetShowHelp(false)
	list.SetShowTitle(false)

	return Model{
		view:      viewInput,
		input:     input,
		inputKeys: newInputKeyMap(),
		ellipsis:  ellipsis,
		list:      list,
		listKeys:  newListKeyMap(),
		help:      help.New(),
	}
}
