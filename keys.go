package main

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Home   key.Binding
	Search key.Binding
	Focus  key.Binding
	Info   key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Home, k.Search, k.Info}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Home, k.Search},
		{k.Focus, k.Info},
	}
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
		key.WithKeys("space"),
		key.WithHelp("space", "get anime info"),
	),
}
