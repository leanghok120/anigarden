package main

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Home   key.Binding
	Search key.Binding
	Focus  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Home, k.Search}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Home, k.Search},
		{k.Focus},
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
}
