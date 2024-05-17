/**
 * Enumeration of the modes Mother can be in.
 */
package model

import "fmt"

type mode int

const (
	prompting mode = iota // default; Mother is processing user inputs alone
	quitting              // Mother is in the process of attempting to cleanly exit
	returning             // child is done and Mother should take over
	handoff               // child is still processing
)

func (m mode) String() string {
	s := ""
	switch m {
	case prompting:
		s = "prompting"
	case quitting:
		s = "quitting"
	case returning:
		s = "returning"
	case handoff:
		s = "handoff"
	default:
		s = fmt.Sprintf("unknown (%d)", m)
	}
	return s
}
