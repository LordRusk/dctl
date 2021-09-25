// package menu gives interfaces for
// integrating menu programs such as
// dmenu into a go program
//
// provides basic dmenu and ui
// implementations
package menu

// this is supposed to be use once like exec.Cmd
//
// a reuseable interface would be more complex.
// though it might be better for some projects,
// it isn't for this one.
type Menu interface {
	// pass is whether or not the prompt should
	// be password protected.
	Prompt(pass bool) error // prompts the user
	Bytes() []byte          // bytes returned, if any
	Ran() bool              // if the user has been prompted
	Error() error           // returns the error, if any.

	Map() Map // can be nil for other use cases
}

// used for keeping track of ui.Menu's
// while keeping the mass amount of
// potential the ui package has
type Map map[string][2]interface{}
