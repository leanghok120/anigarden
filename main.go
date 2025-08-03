package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type anime struct {
	Name string `json:"name"`
	Body string `json:"description"`
}

func (a anime) Title() string {
	return a.Name
}

func (a anime) Description() string {
	return a.Body
}

func (a anime) FilterValue() string {
	return a.Name
}

type model struct {
	list   list.Model
	err    error
	loaded bool
	width  int
	height int
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

// lol "animes"
type (
	errMsg    struct{ err error }
	animesMsg struct{ animes []anime }
)

func fetchHome() tea.Msg {
	res, err := http.Get("https://aniwatch-api-rosy-one.vercel.app/api/v2/hianime/home")
	if err != nil {
		return errMsg{err}
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return errMsg{err}
	}

	var response struct {
		Data struct {
			SpotlightAnimes []anime `json:"spotlightAnimes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return errMsg{err}
	}

	return animesMsg{response.Data.SpotlightAnimes}
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
		l.SetSize(m.width-h, m.height-v)
		m.list = l
		m.loaded = true
	}

	return m, nil
}

func (m model) View() string {
	if m.loaded {
		return docStyle.Render(m.list.View())
	}
	return docStyle.Render("🌱 loading anime list...")
}

func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Println("err: ", err)
		os.Exit(1)
	}
}
