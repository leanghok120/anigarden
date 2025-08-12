package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type page int

const (
	homePage page = iota
	searchPage
	infoPage
	watchlistPage
)

// helper functions
// function to get selected anime and shove it into fetchAnimeInfo or watchAnime or addAnimeToWatchlist
func handleGetAnimeInfo(l list.Model) tea.Cmd {
	if selected, ok := l.SelectedItem().(anime); ok {
		return func() tea.Msg { return fetchAnimeInfo(selected.ID) }
	}
	return nil
}

func handleWatchAnime(l list.Model, animeId string) tea.Cmd {
	if selected, ok := l.SelectedItem().(episode); ok {
		return func() tea.Msg { return watchAnime(selected.ID, animeId) }
	}
	return nil
}

func handleAddToWatchlist(l list.Model) {
	if selected, ok := l.SelectedItem().(anime); ok {
		addAnimeToWatchlist(selected.ID)
	}
}

func handleRemoveFromWatchlist(l list.Model) {
	if selected, ok := l.SelectedItem().(anime); ok {
		removeAnimeFromWatchlist(selected.ID)
	}
}

func setCustomHelp(l *list.Model, page page) {
	switch page {
	case homePage:
		l.AdditionalShortHelpKeys = func() []key.Binding {
			return []key.Binding{keys.Search, keys.Watchlist, keys.AddToWatchlist}
		}
		l.AdditionalFullHelpKeys = func() []key.Binding {
			return []key.Binding{keys.Search, keys.Watchlist, keys.AddToWatchlist, keys.Info}
		}

	case searchPage:
		l.AdditionalShortHelpKeys = func() []key.Binding {
			return []key.Binding{keys.Home, keys.Watchlist, keys.AddToWatchlist}
		}
		l.AdditionalFullHelpKeys = func() []key.Binding {
			return []key.Binding{keys.Home, keys.Watchlist, keys.Focus, keys.Info}
		}

	case infoPage:
		l.AdditionalShortHelpKeys = func() []key.Binding {
			return []key.Binding{keys.Home, keys.Search, keys.Watchlist, keys.Watch}
		}
		l.AdditionalFullHelpKeys = func() []key.Binding {
			return []key.Binding{keys.Home, keys.Search, keys.Watchlist, keys.Watch}
		}

	case watchlistPage:
		l.AdditionalShortHelpKeys = func() []key.Binding {
			return []key.Binding{keys.Home, keys.Search, keys.RemoveFromWatchlist}
		}
		l.AdditionalFullHelpKeys = func() []key.Binding {
			return []key.Binding{keys.Home, keys.Search, keys.RemoveFromWatchlist, keys.Info}
		}
	}
}

// home page
type homeModel struct {
	list    list.Model
	spinner spinner.Model
	err     error
	loaded  bool
	width   int
	height  int
}

func initHomeModel() homeModel {
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return homeModel{spinner: s}
}

func (h homeModel) Update(msg tea.Msg) (homeModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.width = msg.Width
		h.height = msg.Height
		if h.loaded {
			width, height := docStyle.GetFrameSize()
			h.list.SetSize(msg.Width-width, msg.Height-height)
		}

	case animesMsg:
		items := make([]list.Item, len(msg.animes))
		for i, a := range msg.animes {
			items[i] = a
		}
		l := list.New(items, list.NewDefaultDelegate(), 0, 0)
		l.Title = "Home"

		// update list size
		w, v := docStyle.GetFrameSize()
		l.SetSize(h.width-w, h.height-v)

		setCustomHelp(&l, homePage)

		h.list = l
		h.loaded = true

	case tea.KeyMsg:
		if h.list.FilterState() == list.Filtering {
			break
		}
		switch msg.String() {
		case " ", "enter":
			return h, handleGetAnimeInfo(h.list)

		case "a":
			handleAddToWatchlist(h.list)
		}

	case errMsg:
		h.err = msg.err
	}

	if !h.loaded {
		var cmd tea.Cmd
		h.spinner, cmd = h.spinner.Update(msg)
		return h, cmd
	}

	var cmd tea.Cmd
	h.list, cmd = h.list.Update(msg)
	return h, cmd
}

func (h homeModel) View() string {
	if !h.loaded {
		return docStyle.Render(fmt.Sprintf("%s loading anime list...", h.spinner.View()))
	}
	if h.err != nil {
		return docStyle.Render(h.err.Error())
	}
	return docStyle.Render(h.list.View())
}

// search page
type searchModel struct {
	textInput textinput.Model
	list      list.Model
	spinner   spinner.Model
	spinning  bool
	err       error
	loaded    bool
	width     int
	height    int
}

func initSearchModel() searchModel {
	ti := textinput.New()
	ti.Placeholder = "search anime"
	ti.Width = 20
	ti.Cursor.Blink = true

	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return searchModel{textInput: ti, spinner: s}
}

func (s searchModel) Update(msg tea.Msg) (searchModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.textInput.Focused() {
			switch msg.String() {
			case "enter":
				name := s.textInput.Value()
				s.spinner.Tick()
				s.spinning = true
				return s, tea.Batch(s.spinner.Tick, func() tea.Msg { return searchAnime(name) })
			case "esc":
				s.textInput.Blur()
				return s, nil
			}
		}
		if !s.textInput.Focused() && s.list.FilterState() != list.Filtering {
			// when neither input nor filter is focused, allow `t` to focus textInput
			switch msg.String() {
			case "t":
				s.textInput.Focus()
				return s, nil

			case " ", "enter":
				return s, handleGetAnimeInfo(s.list)

			case "a":
				handleAddToWatchlist(s.list)
			}
		}

	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		if s.loaded {
			w, h := docStyle.GetFrameSize()
			s.list.SetSize(s.width-w, s.height-h)
		}

	case searchResultsMsg:
		items := make([]list.Item, len(msg.animes))
		for i, a := range msg.animes {
			items[i] = a
		}
		l := list.New(items, list.NewDefaultDelegate(), 10, 10)
		l.Title = "Results"

		setCustomHelp(&l, searchPage)

		// Get doc padding
		w, v := docStyle.GetFrameSize()
		l.SetSize(s.width-w, s.height-v)

		s.list = l
		s.textInput.Blur()
		s.loaded = true
		s.spinning = false

	case errMsg:
		s.err = msg.err
		return s, nil
	}

	var cmds []tea.Cmd

	var inputCmd tea.Cmd
	var spinnerCmd tea.Cmd
	s.textInput, inputCmd = s.textInput.Update(msg)
	s.spinner, spinnerCmd = s.spinner.Update(msg)

	cmds = append(cmds, inputCmd)
	cmds = append(cmds, spinnerCmd)

	if s.loaded {
		var listCmd tea.Cmd
		s.list, listCmd = s.list.Update(msg)
		cmds = append(cmds, listCmd)
	}

	return s, tea.Batch(cmds...)
}

func (s searchModel) View() string {
	if s.err != nil {
		return docStyle.Render(s.err.Error())
	}
	if s.spinning {
		return docStyle.Render(fmt.Sprintf("%s\n%s searching...", s.textInput.View(), s.spinner.View()))
	}
	if s.loaded {
		return docStyle.Render(fmt.Sprintf("%s\n%s", s.textInput.View(), s.list.View()))
	}
	return docStyle.Render(s.textInput.View())
}

// info page
type infoModel struct {
	id         string
	name       string
	body       string
	genres     []string
	err        error
	leftWidth  int
	rightWidth int
	height     int
	list       list.Model
	spinner    spinner.Model
	spinning   bool
	loaded     bool
}

func initInfoModel(anime anime, width int, height int) infoModel {
	leftWidth := int(float64(width) * 0.4)
	rightWidth := width - leftWidth
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return infoModel{
		id:         anime.ID,
		name:       anime.Name,
		body:       anime.Body,
		genres:     anime.Genres,
		leftWidth:  leftWidth,
		rightWidth: rightWidth,
		height:     height,
		spinner:    s,
		loaded:     false,
	}
}

func (i infoModel) Update(msg tea.Msg) (infoModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == " " || msg.String() == "enter" {
			// Start spinner for launching mpv
			i.spinning = true
			return i, tea.Batch(i.spinner.Tick, handleWatchAnime(i.list, i.id))
		}

	case tea.WindowSizeMsg:
		i.leftWidth = int(float64(msg.Width) * 0.3)
		i.rightWidth = msg.Width - i.leftWidth

	case episodesMsg:
		i.spinning = false
		i.loaded = true

		items := make([]list.Item, len(msg.episodes))
		for i, ep := range msg.episodes {
			items[i] = ep
		}
		l := list.New(items, list.NewDefaultDelegate(), 0, 0)
		l.Title = "Episodes"

		// Update list size
		w, v := docStyle.GetFrameSize()
		l.SetSize(i.rightWidth-w, i.height-v)

		setCustomHelp(&l, infoPage)
		i.list = l

	case errMsg:
		i.err = msg.err
		i.spinning = false
		return i, nil
	}

	var cmds []tea.Cmd

	var spinnerCmd tea.Cmd
	i.spinner, spinnerCmd = i.spinner.Update(msg)
	cmds = append(cmds, spinnerCmd)

	if i.loaded {
		var listCmd tea.Cmd
		i.list, listCmd = i.list.Update(msg)
		cmds = append(cmds, listCmd)
	}

	return i, tea.Batch(cmds...)
}

func (i infoModel) View() string {
	if i.err != nil {
		return docStyle.Render(i.err.Error())
	}

	i.name = lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1).
		Render(i.name)

	genres := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#909090", Dark: "#626262"}).
		Render("Genres: " + strings.Join(i.genres, ","))

	left := lipgloss.NewStyle().
		Width(i.leftWidth).
		MaxWidth(i.leftWidth).
		Render(fmt.Sprintf("%s\n\n%s\n\n%s\n", i.name, genres, i.body))

	gap := lipgloss.NewStyle().Width(4).Render()

	var rightStr string
	right := lipgloss.NewStyle().
		Width(i.rightWidth).
		MaxWidth(i.rightWidth)

	switch {
	case !i.loaded:
		rightStr = right.Render(fmt.Sprintf("%s loading anime episodes...", i.spinner.View()))
	case i.spinning:
		rightStr = right.Render(fmt.Sprintf("%s launching mpv...", i.spinner.View()))
	default:
		rightStr = right.Render(i.list.View())
	}

	return docStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, left, gap, rightStr))
}

// watchlist page
type watchlistModel struct {
	list    list.Model
	spinner spinner.Model
	err     error
	loaded  bool
	width   int
	height  int
}

func initWatchlistModel() watchlistModel {
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return watchlistModel{spinner: s}
}

func (w watchlistModel) Update(msg tea.Msg) (watchlistModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		w.width = msg.Width
		w.height = msg.Height
		if w.loaded {
			width, height := docStyle.GetFrameSize()
			w.list.SetSize(msg.Width-width, msg.Height-height)
		}

	case animesMsg:
		items := make([]list.Item, len(msg.animes))
		for i, a := range msg.animes {
			items[i] = a
		}
		l := list.New(items, list.NewDefaultDelegate(), 0, 0)
		l.Title = "Watchlist"

		// update list size
		wi, vi := docStyle.GetFrameSize()
		l.SetSize(w.width-wi, w.height-vi)

		setCustomHelp(&l, watchlistPage)

		w.list = l
		w.loaded = true

	case tea.KeyMsg:
		if w.list.FilterState() == list.Filtering {
			break
		}
		if msg.String() == " " || msg.String() == "enter" {
			return w, handleGetAnimeInfo(w.list)
		}
		if msg.String() == "r" {
			handleRemoveFromWatchlist(w.list)
			return w, func() tea.Msg { return fetchWatchlist() }
		}

	case errMsg:
		w.err = msg.err
	}

	if !w.loaded {
		var cmd tea.Cmd
		w.spinner, cmd = w.spinner.Update(msg)
		return w, cmd
	}

	var cmd tea.Cmd
	w.list, cmd = w.list.Update(msg)
	return w, cmd
}

func (w watchlistModel) View() string {
	if !w.loaded {
		return docStyle.Render(fmt.Sprintf("%s loading watchlist...", w.spinner.View()))
	}
	if w.err != nil {
		return docStyle.Render(w.err.Error())
	}
	return docStyle.Render(w.list.View())
}
