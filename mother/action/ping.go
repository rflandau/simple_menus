package action

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type Ping struct {
	cobra.Command
}

func (p Ping) Update(m *tea.Model, msg *tea.Msg) (tea.Model, tea.Cmd) {
	return *m, tea.Println("I ping, therefore I am alive")
}

func (p Ping) View(m *tea.Model, msg *tea.Msg) string {
	return "."
}
