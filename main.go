package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

// page enum
type page int

const (
	homePage page = iota
	aboutPage
)

type model struct {
	currPage page
	home     homeModel
	about    aboutModel
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "h":
			m.currPage = homePage
			return m, nil
		case "a":
			m.currPage = aboutPage
			return m, nil
		}
	}

	switch m.currPage {
	case homePage:
		var cmd tea.Cmd
		m.home, cmd = m.home.Update(msg)
		return m, cmd
	case aboutPage:
		var cmd tea.Cmd
		m.about, cmd = m.about.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	switch m.currPage {
	case homePage:
		return m.home.View()
	case aboutPage:
		return m.about.View()
	default:
		return "404 page not found"
	}
}

// pages
type (
	homeModel  struct{}
	aboutModel struct{}
)

// home page
func (h homeModel) Update(msg tea.Msg) (homeModel, tea.Cmd) {
	return h, nil
}

func (h homeModel) View() string {
	return "Home page"
}

// about page
func (a aboutModel) Update(msg tea.Msg) (aboutModel, tea.Cmd) {
	return a, nil
}

func (a aboutModel) View() string {
	return "About page"
}

func main() {
	m := model{home: homeModel{}, about: aboutModel{}}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
