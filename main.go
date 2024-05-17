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
	"simple_menus/model"
	"simple_menus/model/command"

	tea "github.com/charmbracelet/bubbletea"
)

// the top level command
var root *model.Menu

/**
 * generatetree creates the command and menu tree representing the entire CLI.
 * Returns a pointer to the root node
 */
func generateTree() *model.Menu {
	// generate the root of the tree
	r := &model.Menu{
		Name:     "root",
		Parent:   nil,
		Submenus: make(map[string]model.Menu),
		Commands: make(map[string]model.Leaf),
	}

	// generate search submenu
	search := model.Menu{Name: "search", Parent: r, Submenus: nil, Commands: nil}
	// attach it to root
	r.Submenus["search"] = search

	// generate admin submenu
	admin := model.Menu{
		Name:     "admin",
		Parent:   r,
		Submenus: make(map[string]model.Menu),
		Commands: make(map[string]model.Leaf),
	}

	// generate users submenu
	users := model.Menu{Name: "users", Parent: &admin, Submenus: nil, Commands: nil}
	// attach it to admin
	admin.Submenus["users"] = users

	// generate system submenu
	system := model.Menu{Name: "system", Parent: &admin, Submenus: nil, Commands: make(map[string]model.Leaf)}
	system.Commands["status"] = &command.StatusCmd{}

	// attach it to admin
	admin.Submenus["system"] = system

	r.Submenus["admin"] = admin
	return r
}

func init() {
	root = generateTree()
}

func main() {
	var p *tea.Program = tea.NewProgram(model.Initial("simple.log", root))
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
