package recent

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/esferadigital/clima/internal/openmeteo"
	"github.com/esferadigital/clima/internal/store"
)

// ---- keymap ----

type keyMap struct {
	up        key.Binding
	down      key.Binding
	pick      key.Binding
	newSearch key.Binding
	quit      key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.up, k.down, k.pick, k.newSearch, k.quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.up},
		{k.down},
		{k.pick},
		{k.newSearch},
		{k.quit},
	}
}

func newKeyMap() keyMap {
	return keyMap{
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

type RecentCompleteMsg struct {
	Location openmeteo.GeocodingResult
	OK       bool
}

type NewSearchMsg struct{}

// ---- helpers ----

// Implements list.Item interface and wraps location.Location
type recentLocationItem struct {
	openmeteo.GeocodingResult
}

func (i recentLocationItem) FilterValue() string {
	return i.Name
}

func (i recentLocationItem) Title() string {
	return fmt.Sprintf("%s, %s", i.Name, i.Country)
}

func (i recentLocationItem) Description() string {
	return fmt.Sprintf("Lat: %.4f, Lon: %.4f", i.Latitude, i.Longitude)
}

// ---- cmd ----

func getRecentLocationsCmd() tea.Cmd {
	return func() tea.Msg {
		locations, err := store.LoadRecentLocations()
		if err != nil {
			return errorMsg{
				err: err,
			}
		}

		return dataMsg{
			locations: locations,
		}
	}
}

func pickCmd(location openmeteo.GeocodingResult, ok bool) tea.Cmd {
	return func() tea.Msg {
		return RecentCompleteMsg{
			Location: location,
			OK:       ok,
		}
	}
}

func requestNewSearchCmd() tea.Cmd {
	return func() tea.Msg {
		return NewSearchMsg{}
	}
}

// ---- model ----

type view int

const (
	viewList = iota
	viewError
)

type Model struct {
	view view
	list list.Model
	keys keyMap
	help help.Model
}

func (m Model) Init() tea.Cmd {
	return getRecentLocationsCmd()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keys.pick) {
			picked, ok := m.list.SelectedItem().(recentLocationItem)
			if ok {
				return m, pickCmd(picked.GeocodingResult, true)
			}
		}
		if key.Matches(msg, m.keys.newSearch) {
			return m, requestNewSearchCmd()
		}
		if key.Matches(msg, m.keys.quit) {
			return m, tea.Quit
		}
	case dataMsg:
		if len(msg.locations) == 0 {
			return m, pickCmd(openmeteo.GeocodingResult{}, false)
		}
		if len(msg.locations) == 1 {
			return m, pickCmd(msg.locations[0], true)
		}

		items := make([]list.Item, len(msg.locations))
		for i, loc := range msg.locations {
			items[i] = recentLocationItem{loc}
		}
		m.list.SetItems(items)
		return m, nil
	case errorMsg:
		m.view = viewError
		return m, nil
	}

	// Forward messages to sub-components
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	switch m.view {
	case viewList:
		return "\nRecent locations:\n\n" + m.list.View() + "\n" + m.help.View(m.keys)
	case viewError:
		return "\nError occurred"
	default:
		return "Unknown state (recent)"
	}
}

func New() Model {
	listDelegate := list.NewDefaultDelegate()
	listDelegate.ShowDescription = false
	list := list.New([]list.Item{}, listDelegate, 30, 14)
	list.SetShowStatusBar(false)
	list.SetFilteringEnabled(false)
	list.SetShowHelp(false)
	list.SetShowTitle(false)

	return Model{
		view: viewList,
		list: list,
		keys: newKeyMap(),
		help: help.New(),
	}
}
