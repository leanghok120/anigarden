package main

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Home           key.Binding
	Search         key.Binding
	Focus          key.Binding
	Info           key.Binding
	Watchlist      key.Binding
	AddToWatchlist key.Binding
	Watch          key.Binding
}

var keys = keyMap{
	Home: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "home"),
	),
	Search: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "search"),
	),
	Focus: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "focus search bar"),
	),
	Info: key.NewBinding(
		key.WithKeys("space", "enter"),
		key.WithHelp("space/enter", "get anime info"),
	),
	Watchlist: key.NewBinding(
		key.WithKeys("w"),
		key.WithHelp("w", "watchlist"),
	),
	AddToWatchlist: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add to watchlist"),
	),
	Watch: key.NewBinding(
		key.WithKeys("space", "enter"),
		key.WithHelp("space/enter", "watch"),
	),
}
