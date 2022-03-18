// a simple dmenu Menu implementation
package menu

// TODO fix user login with
// dmenu // menu usage

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// errors
var NotInstalled = errors.New("Dmenu not installed")
var AlrRan = errors.New("Already Ran")
var Escaped = errors.New("User made no selection")

type dmenu struct {
	cmd      *exec.Cmd
	bytes    []byte
	ran      bool
	err      error
	passFlag string
	mapp     Map
}

// items seperated by '\n'
//
// mapp can be nil and passFlag can be
// an empty string if not patched
//
// mapp can be nil
func NewDmenu(items string, flags []string, passFlag string, mapp Map) (Menu, error) {
	if _, err := exec.LookPath("dmenu"); err != nil {
		return nil, NotInstalled
	}

	cmd := exec.Command("dmenu", flags[:]...)
	cmd.Stdin = strings.NewReader(items)

	return &dmenu{cmd: cmd, passFlag: passFlag, mapp: mapp}, nil
}

// like NewDmenu but takes []string
// for items rather than a string
// seperated by '\n'
//
// mapp can be nil
func NewSlicedDmenu(items, flags []string, passFlag string, mapp Map) (Menu, error) {
	return NewDmenu(strings.Join(items, "\n"), flags, passFlag, mapp)
}

// prompts the user
//
// returns error if user escapes
func (d *dmenu) Prompt(pass bool) error {
	if d.ran {
		return AlrRan
	}
	if pass {
		d.cmd.Args = append(d.cmd.Args, d.passFlag)
	}

	// dmenu returns an error if the
	// user escapes the selection
	bites, err := d.cmd.CombinedOutput()
	d.bytes = bytes.TrimSpace(bites)
	d.ran = true
	d.err = err

	if err != nil {
		return Escaped
	}
	return nil
}

// returns bytes from users selection
func (d *dmenu) Bytes() []byte { return d.bytes }

// returns whether the menu has been ran
func (d *dmenu) Ran() bool { return d.ran }

// retuns the error, if any
func (d *dmenu) Error() error { return d.err }

// returns the Menu's Map
func (d *dmenu) Map() Map { return d.mapp }
