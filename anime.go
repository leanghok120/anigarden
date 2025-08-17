package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strings"

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

type episode struct {
	ID       string `json:"episodeId"`
	Name     string `json:"title"`
	Number   int    `json:"number"`
	IsFiller bool   `json:"isFiller"`
}

type streamingData struct {
	Data struct {
		Tracks []struct {
			Url  string `json:"file"`
			Lang string `json:"label"`
		} `json:"tracks"`
		Sources struct {
			Url string `json:"file"`
		} `json:"link"`
	} `json:"data"`
}

// lol "animes"
// searchResultsMsg contains anime without a desc
type (
	errMsg           struct{ err error }
	animesMsg        struct{ animes []anime }
	episodesMsg      struct{ episodes []episode }
	searchResultsMsg struct{ animes []anime }
	animeInfoMsg     struct{ anime anime }
	watchlistMsg     struct{ animes []anime }
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

func (e episode) Title() string {
	return fmt.Sprintf("%d. %s", e.Number, e.Name)
}

func (e episode) Description() string {
	if e.IsFiller {
		return "filler"
	}
	return ""
}

func (e episode) FilterValue() string {
	return e.Name
}

// api calls
const (
	url         = "https://aniwatch-api-rosy-one.vercel.app/api/v2/hianime"
	fallbackurl = "https://hianime-api-fallback.onrender.com/api/v1"
)

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

func fetchEpisodes(id string) tea.Msg {
	res, err := http.Get(url + "/anime/" + id + "/episodes")
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
			Episodes []episode `json:"episodes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return errMsg{err}
	}

	return episodesMsg{response.Data.Episodes}
}

// since the api doesn't provide an endpoint that allows
// us to fetch an anime by its id we have to search for animes wih similar id
// and filter the search result until we get the anime we want and repeat till
// we have all the animes in the watchlist
func fetchWatchlist() tea.Msg {
	animeIds := getWatchlist()

	var animesInWatchlist []anime

	// iterate over each anime ID in the watchlist
	for _, animeId := range animeIds {
		res, err := http.Get(fmt.Sprintf("%s/search?q=%s", url, animeId))
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

		// loop through the search results to find the exact match
		for _, foundAnime := range response.Data.Animes {
			if foundAnime.ID == animeId {
				animesInWatchlist = append(animesInWatchlist, foundAnime)
				break
			}
		}
	}

	// Return the final slice of animes wrapped in an animesMsg
	return watchlistMsg{animesInWatchlist}
}

func watchAnime(epId, animeId, lang, client string) tea.Msg {
	res, err := http.Get(fallbackurl + "/stream?id=" + epId + "&server=HD-2&type=" + lang)
	if err != nil {
		return errMsg{err}
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return errMsg{err}
	}

	var response streamingData
	if err := json.Unmarshal(body, &response); err != nil {
		return errMsg{err}
	}

	if client == "mpv" {
		var subFile string
		headers := "Referer: https://vidwish.live/"
		sourceFile := response.Data.Sources.Url

		// get english subtitles
		for _, track := range response.Data.Tracks {
			if track.Lang == "English" {
				subFile = track.Url
			}
		}

		args := []string{"--http-header-fields=" + headers}

		if subFile != "" {
			args = append(args, "--sub-file="+subFile)
		}

		args = append(args, sourceFile)

		mpvCmd := exec.Command("mpv", args...)

		if err := mpvCmd.Run(); err != nil {
			return errMsg{err}
		}
	} else {
		// get episode ID
		parts := strings.Split(epId, "ep=")
		if len(parts) < 2 {
			return errMsg{err}
		}
		epIdNum := parts[1]

		animeUrl := fmt.Sprintf("https://megaplay.buzz/stream/s-2/%s/%s", epIdNum, lang)
		fullUrl := fmt.Sprintf("https://anigarden-player.netlify.app/?iframeLink=%s", animeUrl)

		// Open the browser
		var cmd string
		var args []string
		switch runtime.GOOS {
		case "linux":
			cmd = "xdg-open"
			args = []string{fullUrl}
		case "windows":
			cmd = "rundll32"
			args = []string{"url.dll,FileProtocolHandler", fullUrl}
		case "darwin":
			cmd = "open"
			args = []string{fullUrl}
		default:
			return fmt.Errorf("unsupported platform")
		}

		exec.Command(cmd, args...).Start()
	}

	return fetchEpisodes(animeId) // refetch epiodes after finish watching
}
