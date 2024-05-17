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

// Is there a concern with leaves finishing their processing, calling view, and
// then immediately exiting before a user can actually see the data?

// We could remedy this by requiring users to hit esc or q to return to navigation
// or maintain second model that contains the results of the previous command
// until a new command is issued OR the user enters `clear`
// Finally, we could enter alt mode for leaves and exit on their completion

// Could we model it off of exec? Or are the contexts too seperate?

// tea.Execs?
//	No, these are blocking, but meant for processing and we likely need input
// Could be useful for pinging the server, but there doesn't appear to be a way to get data back
