// ui is a simple abstraction for implementing easily scriptable
// user interfaces in cli/tui enviroments.
//
// implements github.com/lordrusk/dctl/menu
package ui

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"syscall"

	"github.com/pkg/errors"
	"golang.org/x/term"

	"github.com/lordrusk/dctl/menu"
)

type Menu map[interface{}]interface{}

// A handler function.
type HandlerFunc func(Menu) (interface{}, interface{}, error)

// TODO implement menu usage in
// Handler and Screen
type Handler struct {
	sc *bufio.Scanner
	Hf HandlerFunc

	Flags    []string // used for menu
	UseMenu  bool
	PassFlag string
}

// passFlag is the flag passed to menu
// to turn on password input. Can be empty string
func NewHandler(r io.Reader, sf bufio.SplitFunc, flags []string, passFlag string, useMenu bool) *Handler {
	s := bufio.NewScanner(r)
	s.Split(sf)
	return &Handler{sc: s, Flags: flags, PassFlag: passFlag, UseMenu: useMenu}
}

// get the next input as a []byte
func (h *Handler) NextBytes(prompt string) ([]byte, error) {
	if h.UseMenu {
		d, err := menu.NewDmenu("", append(h.Flags, "-p", prompt), h.PassFlag, nil)
		if err != nil {
			return nil, err
		}
		d.Prompt(false)
		return d.Bytes(), err
	}

	if prompt != "" {
		fmt.Println(prompt)
	}
	if !h.sc.Scan() {
		return nil, errors.Wrap(h.sc.Err(), "Nothing to scan!")
	}
	return h.sc.Bytes(), nil
}

// get the next input as a string
func (h *Handler) NextString(prompt string) (string, error) {
	if h.UseMenu {
		d, err := menu.NewDmenu("", append(h.Flags, "-p", prompt), h.PassFlag, nil)
		if err != nil {
			return "", err
		}

		d.Prompt(false)
		return string(d.Bytes()), d.Error()
	}

	if prompt != "" {
		fmt.Println(prompt)
	}
	if !h.sc.Scan() {
		return "", errors.Wrap(h.sc.Err(), "Nothing to scan!")
	}
	return h.sc.Text(), nil
}

// get the next input as a int
func (h *Handler) NextInt(prompt string) (int, error) {
	str, err := h.NextString(prompt)
	if err != nil {
		return 0, err
	}
	if str == "" {
		return 0, errors.New("No input returned!")
	}

	num, err := strconv.Atoi(str)
	if err != nil {
		return 0, errors.Wrap(err, "failed to convert string to int")
	}
	return num, err
}

// menu must be implemented in HandlerFunc
//
// for cleaner code
func (h *Handler) Ask(m Menu) (interface{}, interface{}, error) {
	return h.Hf(m)
}

// basic function to get user input
// allows for password input
func (h *Handler) PassInput(prompt string) ([]byte, error) {
	if h.UseMenu {
		d, err := menu.NewDmenu("", append(h.Flags, "-p", prompt), h.PassFlag, nil)
		if err != nil {
			return nil, err
		}

		d.Prompt(true)
		return d.Bytes(), d.Error()
	}

	fmt.Print(prompt)
	// using golang/x/term.ReadPassword here is the most portable way to do thih.
	bites, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // ui fix
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get next string")
	}
	return bites, nil
}
