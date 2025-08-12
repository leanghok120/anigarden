# 🌸 anigarden

A cozy **TUI anime viewer** written in Go with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and the [HiAnime API](https://github.com/ghoshRitesh12/aniwatch-api).  

Browse, search, and watch anime right in the comfort of your terminal.

## ✨ Features

- **Home View:** See trending or recommended anime right away.  
- **Search View:** Search for your favorite anime.
- **Anime View:** See details about an anime and its episodes.  
- **MPV Integration:** Stream and watch an anime with mpv.

## 📦 Installation

There are 2 ways to install:

### Go

Make sure you have [mpv](https://mpv.io) and [Go](https://go.dev/dl/) installed (version 1.21+ recommended).
Then run:

```sh
go install github.com/leanghok120/anigarden@latest
```

### Releases

Go to the releases and download the binary that fits your machine

## 🚀 Usage

After installing, simply run:

```sh
anigarden
```

## 🗒️ Todos

- [x] home view
- [x] search view
- [x] style the info page
- [x] anime view
- [x] fix kyebinding issues
- [x] watch anime
- [x] add loading spinners
- [x] add favorites/watch list
- [ ] use sakura instead of mpv
- [ ] add the anime poster image to info page (maybe)

## 🤝 Contributing

Contributions are welcome!
If you’d like to help improve anigarden, you can:

1. Fork the repository
2. Commit your changes
3. Open a Pull Request

## 🙏 Acknowledgements

- [Charmbracelet](https://github.com/charmbracelet) for the incredible Bubble Tea TUI framework and other TUI libraries
- [ghoshRitesh12](https://github.com/ghoshRitesh12/aniwatch-api) for the HiAnime API
- [mpv](https://mpv.io) for the powerful media player
