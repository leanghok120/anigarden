package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	list   list.Model
	err    error
	loaded bool
	width  int
	height int
}

func initialModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return fetchHome
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
		if m.loaded {
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.loaded {
			h, v := docStyle.GetFrameSize()
			m.list.SetSize(msg.Width-h, msg.Height-v)
		}

	case animesMsg:
		items := make([]list.Item, len(msg.animes))
		for i, u := range msg.animes {
			items[i] = u
		}
		l := list.New(items, list.NewDefaultDelegate(), 0, 0)
		l.Title = "Home"
		h, v := docStyle.GetFrameSize()
		l.SetSize(m.width-h, m.height-v) // set the size of the list
		m.list = l
		m.loaded = true
	case errMsg:
		m.err = msg.err
	}

	return m, nil
}
