package model

import (
	tea "github.com/charmbracelet/bubbletea"
)

/**
 * A leaf is a command
 * When active, its update, view methods supplant Mother's (the standard model)
 * On completion, leaves MUST send call mother.Return()
 */
type Leaf interface {
	Name() string
	Update(mother *Model, msg tea.Msg) (*Model, tea.Cmd)
	View(mother *Model) string
}

type StatusCmd struct {
	dots string // placeholders while waiting for mother's update to take control
}

func (s *StatusCmd) Name() string {
	return "status"
}

func (s *StatusCmd) Update(mother *Model, msg tea.Msg) (*Model, tea.Cmd) {
	mother.Return()
	return mother, tea.Println("Status: All is well")
}

func (s *StatusCmd) View(mother *Model) string {
	s.dots += "."
	// this will only be visible momentarily
	// anything shown here will be overwritten immediately, hence the tea.Println in Update
	return s.dots
}
