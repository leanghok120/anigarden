package main

import (
	"encoding/json"
	"io"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

type anime struct {
	Name string `json:"name"`
	Body string `json:"description"`
}

// lol "animes"
type (
	errMsg    struct{ err error }
	animesMsg struct{ animes []anime }
)

// list.item implementation
func (a anime) Title() string {
	return a.Name
}

func (a anime) Description() string {
	return a.Body
}

func (a anime) FilterValue() string {
	return a.Name
}

// api calls
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
