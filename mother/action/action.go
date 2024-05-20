package action

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type Actor interface {
	cobra.Command
	Update(m *tea.Model, msg *tea.Msg) (tea.Model, tea.Cmd)
	View(m *tea.Model, msg *tea.Msg) string
}

type Action struct {
	cobra.Command
}

func (a *Action) Update(m *tea.Model, msg *tea.Msg) (tea.Model, tea.Cmd) {
	return *m, nil
}

func (a *Action) View(m *tea.Model, msg *tea.Msg) string {

	return ""
}

/*
type iface interface {
}

type base struct {
	val int
}

type super struct {
	base
	iface
}

func foo(b *base) {

}

func bar(s *super) {
	foo(s)
}

func shoo(s *super) {
	foo(&(s.base))
} */
