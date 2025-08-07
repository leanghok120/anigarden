package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type page int

const (
	homePage page = iota
	searchPage
)

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

		// add custom keybindings to help menu
		l.AdditionalShortHelpKeys = func() []key.Binding {
			return []key.Binding{keys.Home, keys.Search}
		}
		l.AdditionalFullHelpKeys = func() []key.Binding {
			return []key.Binding{keys.Home, keys.Search, keys.Focus}
		}

		h.list = l
		h.loaded = true

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
		switch msg.String() {
		case "enter":
			name := s.textInput.Value()
			return s, func() tea.Msg { return searchAnime(name) }

		case "esc":
			s.textInput.Blur()
			return s, nil

		case "t":
			s.textInput.Focus()
			return s, nil
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

		// add custom keybindings to help menu
		l.AdditionalShortHelpKeys = func() []key.Binding {
			return []key.Binding{keys.Home, keys.Search}
		}
		l.AdditionalFullHelpKeys = func() []key.Binding {
			return []key.Binding{keys.Home, keys.Search, keys.Focus}
		}

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
