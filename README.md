# anigarden

a cozy tui anime viewer written in go with [bubbletea](https://github.com/charmbracelet/bubbletea) and [hianime api](https://github.com/ghoshRitesh12/aniwatch-api)

## todos

- [x] home view
- [x] search view
- [ ] anime view (in progress, got info anime to display anime name)
- [ ] style the info page and maybe add the anime poster image (copy lipgloss style from bubbles)
- [ ] watch anime

## dev notes

### watching

- fetch streaming link from api
- put the Referer header returned from the api in the mpv --http-header-fields (fix 403 error)
- put the .vtt file into mpv --sub-file
- paste the .m3u8 file in
