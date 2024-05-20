/**
 * This is the workhorse of the program.
 * Mother is the Bubble Tea controller, driving Bubble Tea and providing the data
 * representation. It is the default actor and manages passing control to, and
 * retaking control from, child commands invoked by the user. It also contains
 * global data needed to coordinate children and how the program appears.
 * Implements Bubble Tea's model interface
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

	"github.com/spf13/cobra"
)

const (
	//tiCharLimit        = 150
	tiWidth            = 40
	tiValidationString = `^[a-zA-Z\.]+$` // note the anchor wraps
)

// keys that kill the program in Update no matter other states
var killKeys = [...]tea.KeyType{tea.KeyCtrlC, tea.KeyEsc}
var validationRgx = regexp.MustCompile(tiValidationString)

// the data representation of our front-end
type Mother struct {
	ti             textinput.Model
	err            error
	inputErr       error
	log            *log.Logger
	mode           mode
	ss             style.Sheet
	activeCommand  *cobra.Command // nil unless mode == handoff
	builtinActions map[string]func() tea.Cmd

	// Cobra navigation
	root *cobra.Command // graph root; unchanged after initialization
	PWD  *cobra.Command // current command node; should always be a menu
}

var _ tea.Model = Mother{} // compile-time interface check

func textValidator(s string) error {
	if validationRgx.MatchString(s) {
		return nil
	}
	return fmt.Errorf("input contains non-alphabet inputs")
}

/**
 * NewMother generates and returns a controller satisfying Bubble Tea's model
 * interface and operating off the given Cobra tree.
 * @param logpath - path to the file to log to
 * @param root - root node cobra's command tree
 * @param pwd - root if nil, otherwise sets starting submenu
 */
func NewMother(logpath string, root *cobra.Command, pwd *cobra.Command) Mother {
	m := Mother{mode: prompting, root: root, PWD: root}
	if pwd != nil {
		m.PWD = pwd
	}
	// TODO allow a stylesheet to be passed in

	// set up the loggers
	f, err := os.Create(logpath)
	if err != nil {
		panic(err)
	}
	m.log = log.New(f, "", 0)

	/*_, err = tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}*/

	// set up the text input submodel
	m.ti = textinput.New()
	m.ti.Placeholder = "help"
	m.ti.Focus()
	m.ti.Width = tiWidth
	m.ti.Validate = textValidator

	// generate a style sheet
	m.ss.SubmenuText = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAFF"))
	m.ss.CommandText = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFAAAA")).Italic(true)
	m.ss.ErrorText = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF3333"))

	// TODO lock these behind a DEBUG log level
	fmt.Println(m.ss.SubmenuText.Render("set submenu text"))
	fmt.Println(m.ss.CommandText.Render("set command text"))
	fmt.Println(m.ss.ErrorText.Render("set error text"))

	// generate the list of builtin actions
	m.builtinActions = map[string](func() tea.Cmd){
		"..":   m.navParent,
		"help": m.CmdContextHelp,
		"quit": m.quit,
		"exit": m.quit}

	return m
}

/**
 * Called by leaves to return handling to the standard/navigation model
 */
func (m *Mother) Return() {
	m.mode = returning
}

//#region tea interface

func (m Mother) Init() tea.Cmd {
	return textinput.Blink
}

/* Inputs are handled in two places:
 * Keystroke commands (ex: F1, CTRL+C) are handled here.
 * Input commands (ex: 'help', 'quit', <command>) are handled in processInput(),
 * be they built-in or local commands */
func (m Mother) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

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

	// check if child action is done
	if m.mode == returning {
		m.log.Println("Returning from command...")
		// ensure we are in an active command
		if m.activeCommand == nil {
			panic("return mode but no active command")
		}
		m.activeCommand = nil
		m.mode = prompting
	}
	// allow child action to retain control if it exists
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
		if msg.Type == tea.KeyF1 { // help
			return m, m.CmdContextHelp()
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

func (m Mother) View() string {
	// allow child command to retain control if it exists
	if m.activeCommand != nil {
		return m.activeCommand.View(&m)
	}

	// if there was a fatal error, print it out and return
	if m.err != nil {
		return fmt.Sprintf("\nErr: %v\n\n", m.err)
	}

	s := strings.Builder{}
	s.WriteString(m.PWD.Name() + " " + m.ti.View()) // prompt
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

/**
 * processInput consumes and clears the text on the prompt, determining what action to take and modifying the model accordingly.
 * If we change directory, prints the new pwd above
 */
func processInput(m *Mother) tea.Cmd {
	var given string = m.ti.Value()
	m.ti.Validate(given)
	if m.ti.Err != nil {
		return nil
	}
	m.ti.Reset() // empty out the input
	// check for a builtin command
	builtinFunc, ok := m.builtinActions[given]
	if ok {
		return builtinFunc()
	}
	// if we do not find a built in, test for a valid child
	var child *cobra.Command = nil
	for _, c := range m.PWD.Commands() {
		m.log.Printf("Given '%s' =?= child '%s'", given, c.Name()) // DEBUG

		if c.Name() == given { // match
			m.log.Printf(" | true\n", given, c.Name()) // DEBUG
			child = c
			break
		}
		m.log.Printf("\n", given, c.Name()) // DEBUG
	}

	// check if we found a match
	if child == nil {
		// user request unhandlable
		m.inputErr = fmt.Errorf("%s has no child '%s'", m.PWD.Name(), given)
		return nil
	}

	// split on action or nav
	if isAction(child) {
		// hand off control to child
		m.log.Printf("Found local command %v\n", child.Name())
		m.mode = handoff
		// TODO each time a command is call, it should be instantiated fresh so
		//	old data does not garble the new call
		m.activeCommand = child
		return nil
	} else { // nav
		// navigate to child
		m.PWD = child
		return m.CmdPWD()
	}
}

/** Returns a tea.Println Cmd containing the context help for the command pointed to by PWD */
func (m *Mother) CmdContextHelp() tea.Cmd {
	return tea.Println("Help for " + m.PWD.Name())
}

/** Returns a tea.Println Cmd containing the path to the pwd */
func (m *Mother) CmdPWD() tea.Cmd {
	return tea.Println(m.PWD.CommandPath())
}

/* Using the current menu, navigate up one level */
func (m *Mother) navParent() tea.Cmd {
	if m.PWD == m.root { // if we are at root, do nothing
		return nil
	}
	// otherwise, step upward
	m.PWD = m.PWD.Parent()
	return m.CmdPWD()
}

/* Quit the program */
func (m *Mother) quit() tea.Cmd {
	return tea.Sequence(tea.Println("Bye"), tea.Quit)
}

// #region helper subroutines

/**
 * Given a cobra.Command, returns whether it is an Action (and thus its .Run()
 * can be called) or a Nav (and its .Run() would be redundant if we are already
 * in interactive mode)
 */
func isAction(cmd *cobra.Command) bool {
	if cmd == nil { // sanity check
		panic("cmd cannot be nil!")
	}
	if cmd.ContainsGroup("action") {
		return true
	} else if cmd.ContainsGroup("nav") {
		return false
	} else { // sanity check
		panic("cmd '" + cmd.Name() + "' is neither a nav nor an action!")
	}
}

//#endregion
