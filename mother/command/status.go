/**
 * Determines and reports on the current system status.
 * As with all Commands, implements the leaf interface.
 */

package command

import (
	. "simple_menus/mother"

	tea "github.com/charmbracelet/bubbletea"
)

var _ Leaf = &StatusCmd{} // compile-time interface fulfillment check

type StatusCmd struct {
	dots string // placeholders while waiting for mother's update to take control
}

func (s *StatusCmd) Name() string {
	return "status"
}

func (s *StatusCmd) Update(mother *Mother, msg tea.Msg) (*Mother, tea.Cmd) {
	mother.Return()
	return mother, tea.Println("Status: All is well")
}

func (s *StatusCmd) View(mother *Mother) string {
	s.dots += "."
	// this will only be visible momentarily
	// anything shown here will be overwritten immediately, hence the tea.Println in Update
	return s.dots
}
