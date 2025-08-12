package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	currPage  page
	home      homeModel
	search    searchModel
	info      infoModel
	watchlist watchlistModel
	win       tea.WindowSizeMsg
}

func initialModel() model {
	return model{currPage: homePage, home: initHomeModel(), search: initSearchModel(), watchlist: initWatchlistModel()}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(fetchHome, m.home.spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.win = msg

	case animeInfoMsg:
		m.info = initInfoModel(msg.anime, m.win.Width, m.win.Height)
		m.currPage = infoPage
		return m, tea.Batch(m.info.spinner.Tick, func() tea.Msg { return fetchEpisodes(msg.anime.ID) })

	case tea.KeyMsg:
		// if in filtering or textinput focus state, avoid quiting, switch pages...
		if m.home.list.FilterState() == list.Filtering || m.search.list.FilterState() == list.Filtering || m.search.textInput.Focused() {
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
			m.search.textInput.Focus()
			return m, func() tea.Msg { return m.win } // send tea.WindowSizeMsg to search model

		case "w":
			m.currPage = watchlistPage

			// send tea.WindowSizeMsg to watchlist model
			return m, tea.Batch(fetchWatchlist, func() tea.Msg { return m.win }, m.watchlist.spinner.Tick)
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

	case infoPage:
		var cmd tea.Cmd
		m.info, cmd = m.info.Update(msg)
		return m, cmd

	case watchlistPage:
		var cmd tea.Cmd
		m.watchlist, cmd = m.watchlist.Update(msg)
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
	case infoPage:
		return m.info.View()
	case watchlistPage:
		return m.watchlist.View()
	default:
		return "404 not found"
	}
}
