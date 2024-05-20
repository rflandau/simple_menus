package mother

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
