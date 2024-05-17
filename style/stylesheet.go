package style

import "github.com/charmbracelet/lipgloss"

type Sheet struct {
	SubmenuText lipgloss.Style
	CommandText lipgloss.Style
	ErrorText   lipgloss.Style
}

// TODO generate different stylesheet functions that return a prebuild StyleSheet
//func Monokai() Sheet {}
