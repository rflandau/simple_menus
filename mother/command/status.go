/**
 * Determines and reports on the current system status.
 * As with all Commands, implements the leaf interface.
 */

package command

import (
	. "simple_menus/mother"

	tea "github.com/charmbracelet/bubbletea"
)

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
