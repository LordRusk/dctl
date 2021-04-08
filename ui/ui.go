// ui is a simple abstraction for implementing easily scriptable
// user interfaces in cli/tui enviroments.
package ui

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"syscall"

	"github.com/pkg/errors"
	"golang.org/x/term"
)

type Menu map[interface{}]interface{}

// A handler function.
type HandlerFunc func(Menu) (interface{}, interface{}, error)

type Handler struct {
	sc *bufio.Scanner
	Hf HandlerFunc
}

func NewHandler(r io.Reader, sf bufio.SplitFunc) *Handler {
	s := bufio.NewScanner(r)
	s.Split(sf)
	return &Handler{sc: s}
}

// get the next input as a []byte
func (h *Handler) NextBytes() ([]byte, error) {
	if !h.sc.Scan() {
		return nil, errors.New("Nothing to scan!")
	}
	return h.sc.Bytes(), nil
}

// get the next input as a string
func (h *Handler) NextString() (string, error) {
	if !h.sc.Scan() {
		return "", errors.New("Nothing to scan!")
	}
	return h.sc.Text(), nil
}

func (h *Handler) NextInt() (int, error) {
	str, err := h.NextString()
	if err != nil {
		return 0, err
	}
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0, errors.Wrap(err, "failed to convert string to int")
	}
	return num, err
}

// bad name, I know.
type Screen struct {
	*Handler
}

func New(h *Handler) *Screen {
	return &Screen{Handler: h}
}

// most basic use of the ui package
func (s *Screen) Ask(m Menu) (interface{}, interface{}, error) {
	return s.Hf(m)
}

// basic function to get user input
func (s *Screen) Input(pass bool) (string, error) {
	var str string
	var err error
	if pass {
		// using golang/x/term.ReadPassword here is the most portable way to do this.
		var bites []byte
		bites, err = term.ReadPassword(int(syscall.Stdin))
		str = string(bites)
		fmt.Println() // ui fix
	} else {
		str, err = s.NextString()
	}
	if err != nil {
		return "", errors.Wrap(err, "Failed to get next string")
	}
	return str, nil
}
