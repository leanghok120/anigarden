package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	initDB()
	defer db.Close()

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatalf("failed to run anigarden: %v\n", err)
	}
}
