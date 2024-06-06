package program

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type ProgramState struct {
	ExitState     bool
	Username      string
	DataToCollect []string
}

func (p *ProgramState) ExitIfRequested(tprogram *tea.Program) {
	if p.ExitState {
		if err := tprogram.ReleaseTerminal(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Good Bye!")
		os.Exit(1)
	}
}
