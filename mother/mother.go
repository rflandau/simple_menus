/**
 * This is the workhorse of the program.
 * Mother is the Bubble Tea controller, driving Bubble Tea and providing the data
 * representation. It is the default actor and manages passing control to, and
 * retaking control from, child commands invoked by the user. It also contains
 * global data needed to coordinate children and how the program appears.
 *
 */
package mother

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"simple_menus/message"
	"simple_menus/style"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	tiCharLimit        = 150
	tiWidth            = 20
	tiValidationString = `^[a-zA-Z\.]+$` // note the anchor wraps
)

// keys that kill the program in Update no matter other states
var killKeys = [...]tea.KeyType{tea.KeyCtrlC, tea.KeyEsc}
var validationRgx = regexp.MustCompile(tiValidationString)

// the data representation of our front-end
type Model struct {
	ti            textinput.Model
	err           error
	inputErr      error
	curMenu       *Menu
	log           *log.Logger
	mode          mode
	ss            style.Sheet
	activeCommand Leaf
}

func textValidator(s string) error {
	if validationRgx.MatchString(s) {
		return nil
	}
	return fmt.Errorf("input contains non-alphabet inputs")
}

// TODO allow a stylesheet to be passed in
func Initial(logpath string, root *Menu) Model {
	m := Model{}

	// set up the loggers
	f, err := os.Create(logpath)
	if err != nil {
		panic(err)
	}
	m.log = log.New(f, "", 0)
	// TODO close log files
	_, err = tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}

	// set up the text input submodel
	m.ti = textinput.New()
	m.ti.Placeholder = ""
	m.ti.Focus()
	m.ti.CharLimit = tiCharLimit
	m.ti.Width = tiWidth
	m.ti.Validate = textValidator

	// start on the root node
	m.curMenu = root

	// generate a style sheet
	m.ss.SubmenuText = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAFF"))
	m.ss.CommandText = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFAAAA")).Italic(true)
	m.ss.ErrorText = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF3333"))

	// TODO remove these debug calls
	fmt.Println(m.ss.SubmenuText.Render("set submenu text"))
	fmt.Println(m.ss.CommandText.Render("set command text"))
	fmt.Println(m.ss.ErrorText.Render("set error text"))

	// enter prompt mode
	m.mode = prompting

	return m
}

/**
 * Called by leaves to return handling to the standard/navigation model
 */
func (m *Model) Return() {
	m.mode = returning
}

//#region tea interface

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

/* Inputs are handled in two places:
 * Keystroke commands (ex: F1, CTRL+C) are handled here.
 * Input commands (ex: 'help', 'quit', <command>) are handled in processInput(),
 * be they built-in or local commands */
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//m.log.Printf("Received message %#v\n", msg)

	// always handle kill keys
	keyMsg, isKeyMsg := msg.(tea.KeyMsg)
	if isKeyMsg {
		for _, v := range killKeys {
			if keyMsg.Type == v {
				m.mode = quitting
				return m, tea.Batch(tea.Quit, tea.Println("Bye"))
			}
		}
	}

	// check if child command is done
	if m.mode == returning {
		m.log.Println("Returning from command...")
		// ensure we are in an active command
		if m.activeCommand == nil {
			panic("return message but no active command")
		}
		m.activeCommand = nil
		m.mode = prompting
	}
	// allow child command to retain control if it exists
	if m.activeCommand != nil && m.mode == handoff {
		m.log.Printf("Handing off Update control to active command %s\n", m.activeCommand.Name())
		return m.activeCommand.Update(&m, msg)
	} else if m.activeCommand != nil || m.mode == handoff {
		// if one but not the other, something has gone wrong
		panic(fmt.Sprintf("active command (%s) and mode (%s) are inconsistent", m.activeCommand.Name(), m.mode))
	}

	// normal handling
	switch msg := msg.(type) {
	case message.Err:
		m.err = msg
		return m, tea.Sequence(tea.Println("Bye"), tea.Quit)
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC { // ! kill
			// received kill key
			m.mode = quitting
			return m, tea.Batch(tea.Quit, tea.Println("Bye"))
		}
		if msg.Type == tea.KeyF1 || msg.Type == tea.KeyDelete { // help
			return m, m.curMenu.Children(m.ss)
		}
		if msg.Type == tea.KeyEnter { // submit
			m.log.Printf("User hit enter, parsing data '%v'\n", m.ti.Value())
			// on enter, clear any error text, process the data in the text
			// input, and manipulate model accordingly

			m.inputErr = nil
			cmd := processInput(&m)
			if m.inputErr != nil {
				m.log.Printf("%v\n", m.inputErr)
			}
			return m, cmd
		}
	}

	var cmd tea.Cmd
	m.ti, cmd = m.ti.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	// allow child command to retain control if it exists
	if m.activeCommand != nil {
		return m.activeCommand.View(&m)
	}

	// if there was a fatal error, print it out and return
	if m.err != nil {
		return fmt.Sprintf("\nErr: %v\n\n", m.err)
	}

	s := strings.Builder{}
	s.WriteString(m.curMenu.Name + " " + m.ti.View()) // prompt
	if m.ti.Err != nil {
		// write out previous error and clear it
		s.WriteString("\n")
		s.WriteString(m.ss.ErrorText.Render(m.ti.Err.Error()))
		m.ti.Err = nil
		// this will be cleared from view automagically on next key input
	}
	return s.String()
}

//#endregion

type builtinFunc func(m *Model) tea.Cmd

/**
 * Built in actions avaialble to all menus
 * command -> operator function
 */
var builtin = map[string]builtinFunc{
	"..":   navParent, // walk up a level
	"help": help,
	"quit": quit,
	"exit": quit,
}

/**
 * processInput consumes and clears the text on the prompt, determining what action to take and modifying the model accordingly.
 * If we change directory, prints the new pwd above
 */
func processInput(m *Model) tea.Cmd {
	var s string = m.ti.Value()
	m.ti.Validate(s)
	if m.ti.Err != nil {
		return nil
	}
	m.ti.Reset() // empty out the input
	// check for a builtin command
	builtinFunc, ok := builtin[s]
	if ok {
		return builtinFunc(m)
	}
	// if we do not find a built in, test for a submenu
	m.log.Printf("Parsing for submenus (from: %+v)\n", m.curMenu.Submenus)
	if submenu, ok := m.curMenu.Submenus[strings.ToLower(s)]; ok { // submenu
		m.curMenu = &submenu
		// TODO as well as printing pwd, also print the current contents of ti's buffer (before resetting)
		// 	this should cause a terminal like appearance
		return tea.Println(m.pwd())
	}
	// test for command
	m.log.Printf("Parsing for commands (from: %+v)\n", m.curMenu.Commands)
	if command, ok := m.curMenu.Commands[strings.ToLower(s)]; ok { // command
		// When a command is issued, set the model's active command
		// While a command is set, the model will call its Update and View functions
		// On completion, a command will set m.mode to returning

		// TODO upgrade Help() to support command context by checking m.activeCommand
		// Help() can automatically check for an active command and show the command's help field instead.

		m.log.Printf("Found local command %v\n", command.Name())
		m.mode = handoff
		// TODO each time a command is call, it should be instantiated fresh so
		//	old data does not garble the new call
		m.activeCommand = command
		return nil
	}

	// no child found
	m.inputErr = fmt.Errorf("%s has no child '%+s'", m.curMenu.Name, s)
	// TODO put this inputerr out via View
	return nil
}

/**
 * Returns present working directory.
 */
func (m *Model) pwd() string {
	return m.curMenu.Path()
}

/* Using the current menu, navigate up one level */
func navParent(m *Model) tea.Cmd {
	if m.curMenu.Parent == nil { // if we are at root, do nothing
		return nil
	}
	// otherwise, set upward
	m.curMenu = m.curMenu.Parent
	return tea.Println(m.pwd())
}

/* Print context help for the current menu */
func help(m *Model) tea.Cmd {
	return m.curMenu.Children(m.ss)
}

/* Quit the program */
func quit(m *Model) tea.Cmd {
	return tea.Sequence(tea.Println("Bye"), tea.Quit)
}