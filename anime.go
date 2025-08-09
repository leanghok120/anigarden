package main

import (
	"encoding/json"
	"io"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

// api returns a desc in the home route but doesn't return
// a desc in the search route
type anime struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Body   string   `json:"description"`
	Genres []string `json:"genres"`
}

// lol "animes"
// searchResultsMsg contains anime without a desc
type (
	errMsg           struct{ err error }
	animesMsg        struct{ animes []anime }
	searchResultsMsg struct{ animes []anime }
	animeInfoMsg     struct{ anime anime }
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
const url = "https://aniwatch-api-rosy-one.vercel.app/api/v2/hianime"

func fetchHome() tea.Msg {
	res, err := http.Get(url + "/home")
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

func searchAnime(name string) tea.Msg {
	res, err := http.Get(url + "/search?q=" + name)
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
			Animes []anime `json:"animes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return errMsg{err}
	}

	return searchResultsMsg{response.Data.Animes}
}

func fetchAnimeInfo(id string) tea.Msg {
	res, err := http.Get(url + "/qtip/" + id)
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
			Anime anime `json:"anime"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return errMsg{err}
	}

	return animeInfoMsg{response.Data.Anime}
}
