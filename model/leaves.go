package model

import (
	tea "github.com/charmbracelet/bubbletea"
)

// a leaf is a command
// When active, its update and view methods supplant that of its mother model
// On completion, leaves MUST send a returnMsg
type Leaf interface {
	Name() string
	Update(m *Model, msg tea.Msg) (*Model, tea.Cmd)
	View(m *Model) string
}

type StatusCmd struct{}

func (s StatusCmd) Name() string {
	return "status"
}

func (s StatusCmd) Update(m *Model, msg tea.Msg) (*Model, tea.Cmd) {
	m.Return()
	return m, tea.Println("Status: All is well")
}

func (s StatusCmd) View(m *Model) string {
	// anything shown here will be overwritten immediately, hence the tea.Println in Update
	return "Other all is well"
}
