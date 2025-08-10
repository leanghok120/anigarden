package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type page int

const (
	homePage page = iota
	searchPage
	infoPage
)

// helper functions
// function to get selected anime and shove it into fetchAnimeInfo
func handleAnimeSelection(l list.Model) tea.Cmd {
	if selected, ok := l.SelectedItem().(anime); ok {
		return func() tea.Msg { return fetchAnimeInfo(selected.ID) }
	}
	return nil
}

func setCustomHelp(l *list.Model) {
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{keys.Home, keys.Search, keys.Info}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{keys.Home, keys.Search, keys.Focus, keys.Info}
	}
}

// home page
type homeModel struct {
	list   list.Model
	err    error
	loaded bool
	width  int
	height int
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

		setCustomHelp(&l)

		h.list = l
		h.loaded = true

	case tea.KeyMsg:
		if h.list.FilterState() == list.Filtering {
			break
		}
		if msg.String() == " " || msg.String() == "enter" {
			return h, handleAnimeSelection(h.list)
		}

	case errMsg:
		h.err = msg.err
	}

	if h.loaded {
		var cmd tea.Cmd
		h.list, cmd = h.list.Update(msg)
		return h, cmd
	}

	return h, nil
}

func (h homeModel) View() string {
	if !h.loaded {
		return docStyle.Render("🌱 loading anime list...")
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
	return searchModel{textInput: ti}
}

func (s searchModel) Update(msg tea.Msg) (searchModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.textInput.Focused() {
			switch msg.String() {
			case "enter":
				name := s.textInput.Value()
				return s, func() tea.Msg { return searchAnime(name) }
			case "esc":
				s.textInput.Blur()
				return s, nil
			}
		}
		if !s.textInput.Focused() && s.list.FilterState() != list.Filtering {
			// when neither input nor filter is focused, allow `t` to focus textInput
			if msg.String() == "t" {
				s.textInput.Focus()
				return s, nil
			}
			if msg.String() == " " || msg.String() == "enter" {
				return s, handleAnimeSelection(s.list)
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

		setCustomHelp(&l)

		// Get doc padding
		w, v := docStyle.GetFrameSize()
		l.SetSize(s.width-w, s.height-v)

		s.list = l
		s.textInput.Blur()
		s.loaded = true

	case errMsg:
		s.err = msg.err
		return s, nil
	}

	var cmds []tea.Cmd

	var inputCmd tea.Cmd
	s.textInput, inputCmd = s.textInput.Update(msg)
	cmds = append(cmds, inputCmd)

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
	if s.loaded {
		return docStyle.Render(s.textInput.View() + "\n" + s.list.View())
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
	loaded     bool
}

func initInfoModel(anime anime, width int, height int) infoModel {
	leftWidth := int(float64(width) * 0.4)
	rightWidth := width - leftWidth

	return infoModel{
		id:         anime.ID,
		name:       anime.Name,
		body:       anime.Body,
		genres:     anime.Genres,
		leftWidth:  leftWidth,
		rightWidth: rightWidth,
		height:     height,
		loaded:     false,
	}
}

func (i infoModel) Update(msg tea.Msg) (infoModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		i.leftWidth = int(float64(msg.Width) * 0.4)
		i.rightWidth = msg.Width - i.leftWidth

	case episodesMsg:
		items := make([]list.Item, len(msg.episodes))
		for i, a := range msg.episodes {
			items[i] = a
		}
		l := list.New(items, list.NewDefaultDelegate(), 0, 0)
		l.Title = "Episodes"

		// update list size
		w, v := docStyle.GetFrameSize()
		l.SetSize(i.rightWidth-w, i.height-v)

		// setCustomHelp(&l)

		i.list = l
		i.loaded = true
	}

	var cmd tea.Cmd
	i.list, cmd = i.list.Update(msg)
	return i, cmd
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
		Render(fmt.Sprintf(
			"%s\n\n%s\n\n%s\n",
			i.name,
			genres,
			i.body,
		))

	gap := lipgloss.NewStyle().Width(4).Render()

	var rightStr string
	right := lipgloss.NewStyle().
		Width(i.rightWidth).
		MaxWidth(i.rightWidth)

	if !i.loaded {
		rightStr = right.Render("🌱 loading anime episodes...")
	} else {
		rightStr = right.Render(i.list.View())
	}

	return docStyle.Render(lipgloss.JoinHorizontal(lipgloss.Left, left, gap, rightStr))
}
