package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rdlucas2/doomscroller/internal/engine"

	// Register all story content via init().
	_ "github.com/rdlucas2/doomscroller/content/stories"
)

func main() {
	e := engine.New()
	p := tea.NewProgram(e, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}
