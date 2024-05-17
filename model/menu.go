/**
 * A Menu is a single node within the command tree.
 * It cannot be invoked as it is not a command. However, it may contain commands (leaves).
 * As the command tree is static, menus can lazy-compile `compiled...` fields
 * the first time they are requested. */
package model

import (
	"simple_menus/style"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// menu represents a single command or subcommand within the tree
type Menu struct {
	Name             string          // my name
	Parent           *Menu           // ptr to my parent, for navigation
	Submenus         map[string]Menu // menus beneath me
	Commands         map[string]Leaf // commands callable from me
	compiledChildren string          // sorted list of menus+commands, generated on first call to Children
	compiledPath     string          // path to this command, generated on first call to Path()
}

/* Returns a tea.Cmd to print out the children of the given menu.
 * Lazy-generates the list of children, so it is only requires generation once,
 * on first visit */
func (menu *Menu) Children(s style.Sheet) tea.Cmd {
	if menu.compiledChildren == "" {
		// generate the list of available menus and commands
		allChoices := make([]string, len(menu.Submenus)+len(menu.Commands))

		var i int = 0
		for k := range menu.Submenus {
			allChoices[i] = s.SubmenuText.Render(k)
			i++
		}
		for k := range menu.Commands {
			allChoices[i] = s.CommandText.Render(k)
			i++
		}

		// ! going forward, use the slices package from 1.21+
		sort.Strings(allChoices)
		menu.compiledChildren = strings.Join(allChoices, "\n")
	}
	return tea.Println(menu.compiledChildren)
}

func (m *Menu) Path() string {
	if m.compiledPath == "" {
		// Climbs the tree from current menu, tracking each ancestor
		// Because we can be at an arbitrary depth, we need to append to an array and then reverse it
		path := []string{}
		cur := m
		for cur != nil {
			// append current menu
			path = append(path, cur.Name)
			// climb
			cur = cur.Parent
		}
		// reverse the path to reflect downward trajectory
		for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
			path[i], path[j] = path[j], path[i]
		}
		m.compiledPath = strings.Join(path, "/")
	}
	return m.compiledPath
}
