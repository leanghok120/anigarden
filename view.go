package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type page int

const (
	homePage page = iota
	searchPage
)

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
	return docStyle.Render(h.list.View())
}

type searchModel struct{}

func (s searchModel) Update(msg tea.Msg) (searchModel, tea.Cmd) {
	return s, nil
}

func (s searchModel) View() string {
	return docStyle.Render("search page")
}
