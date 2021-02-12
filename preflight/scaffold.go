package preflight

import "os"

type scaffold struct {
	OSExit func(int)
}

func init() {
	Restore()
}

// Scaffold a set of functions that may be replaced for testing
var Scaffold *scaffold

// Restore restores the scaffold to its default state
func Restore() {
	Scaffold = &scaffold{
		OSExit: os.Exit,
	}
}
