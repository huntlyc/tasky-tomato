package main

import (
	"fmt"
	"log"
	"os"

	"charm.land/bubbletea/v2"

	"github.com/huntlyc/tasky-tomato/internal/db"
	"github.com/huntlyc/tasky-tomato/internal/ui"
)

func main() {
	database, err := db.Open("./tasky.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	if err := db.Init(database); err != nil {
		log.Fatal(err)
	}

	m := ui.NewModel(database)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
