package main

import "github.com/charmbracelet/lipgloss"

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func (m model) View() string {
	if m.loaded {
		return docStyle.Render(m.list.View())
	}
	return docStyle.Render("🌱 loading anime list...")
}
