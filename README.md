# anigarden

a cozy tui anime viewer written in go with [bubbletea](https://github.com/charmbracelet/bubbletea) and [hianime api](https://github.com/ghoshRitesh12/aniwatch-api)

## todos

- [x] home view
- [x] search view
- [x] style the info page
- [x] anime view
- [ ] fix kyebinding issues
- [ ] watch anime
- [ ] add the anime poster image to info page (maybe)

## dev notes

### keybind issues

- pressing t in search bar does not work
- pressing space in search bar and filter does not work
- pressing t while filtering focuses search bar
- pressing esc or enter in filter in search page doesnt work

### watching

- fetch streaming link from api
- put the Referer header returned from the api in the mpv --http-header-fields (fix 403 error)
- put the .vtt file into mpv --sub-file
- paste the .m3u8 file in
