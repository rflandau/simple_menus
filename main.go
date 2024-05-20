/**
 * Another example Bubble Tea driver for navigating multiple menus via Bubble's textinput prompt.
 * Root
 * |
 * |-- Admin
 * |	|-- Users
 * |	|-- System
 * |	|	 |-- Status
 * |	|	 |-- Hardware
 * |-- Search
 *
 */
package main

import (
	"simple_menus/mother"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

/**
 * generatetree creates the command and menu tree representing the entire CLI.
 * Returns a pointer to the root node
 */
func generateTree() *cobra.Command {

	r := &cobra.Command{}

	/*
		// generate the root of the tree
		r := &mother.Menu{
			Name:     "root",
			Parent:   nil,
			Submenus: make(map[string]mother.Menu),
			Commands: make(map[string]mother.Leaf),
		}

		// generate search submenu
		search := mother.Menu{Name: "search", Parent: r, Submenus: nil, Commands: nil}
		// attach it to root
		r.Submenus["search"] = search

		// generate admin submenu
		admin := mother.Menu{
			Name:     "admin",
			Parent:   r,
			Submenus: make(map[string]mother.Menu),
			Commands: make(map[string]mother.Leaf),
		}

		// generate users submenu
		users := mother.Menu{Name: "users", Parent: &admin, Submenus: nil, Commands: nil}
		// attach it to admin
		admin.Submenus["users"] = users

		// generate system submenu
		system := mother.Menu{Name: "system", Parent: &admin, Submenus: nil, Commands: make(map[string]mother.Leaf)}
		system.Commands["status"] = &command.StatusCmd{}

		// attach it to admin
		admin.Submenus["system"] = system

		r.Submenus["admin"] = admin
	*/
	return r
}

func main() {
	var p *tea.Program = tea.NewProgram(mother.NewMother("simple.log", generateTree(), nil))
	_, err := p.Run()
	if err != nil {
		panic(err)
	}

}

/**
 * Workhorse for String; recursively prints the current comm at the given depth level, allowing
 * String() to iterate cleanly.
 */
/*func (c *menu) StringDepth(prevIndent string) (full string, curIndent string) {
	s := strings.Builder{}
	s.WriteString(c.name)
	s.WriteString("\n")
	// for each child, generate a tree of it and all grandchildren
	for k, v := range c.children {
		s.WriteString(prevIndent)
		s.WriteRune('|')
		s.WriteString("-- " + k)
		s.WriteString(v.StringDepth(depth))
	}
}

func (c *menu) String() string {
	// generate the first indent
	indent := strings.Builder{}
	var i uint = 0
	for ; i < indentWidth; i++ {
		indent.WriteRune(' ')
	}
	ret, _ := c.StringDepth(indent.String())
	return ret
}*/
