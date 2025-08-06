package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	currPage page
	home     homeModel
	search   searchModel
}

func initialModel() model {
	return model{currPage: homePage, home: homeModel{}, search: searchModel{}}
}

func (m model) Init() tea.Cmd {
	return fetchHome
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// if in filtering state, avoid quiting, switch pages...
		if m.currPage == homePage && m.home.list.FilterState() == list.Filtering {
			break
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "h":
			m.currPage = homePage
			return m, nil

		case "s":
			m.currPage = searchPage
			return m, nil
		}
	}

	switch m.currPage {
	case homePage:
		var cmd tea.Cmd
		m.home, cmd = m.home.Update(msg)
		return m, cmd
	case searchPage:
		var cmd tea.Cmd
		m.search, cmd = m.search.Update(msg)
		return m, cmd
	}

	return m, nil
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func (m model) View() string {
	switch m.currPage {
	case homePage:
		return m.home.View()
	case searchPage:
		return m.search.View()
	default:
		return "404 not found"
	}
}
