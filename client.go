package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/diamondburned/arikawa/v2/state"
	"github.com/lordrusk/dctl/logger"
	"github.com/lordrusk/dctl/ui"
)

type client struct {
	*state.State
	*ui.Screen
	*logger.Logger
}

func setupDefaultHandlerFunc(h *ui.Handler) {
	// key: string, value: interface{}
	h.Hf = func(m ui.Menu) (interface{}, interface{}, error) { // ui.HandlerFunc
		for {
			var b strings.Builder
			for k, v := range m {
				b.WriteString(fmt.Sprintf("[%s] %s\n", k, v))
			}
			fmt.Print(b.String())
			str, err := h.NextString()
			if err != nil {
				return nil, nil, errors.Wrap(err, "Failed to get next string")
			}
			if m[str] == nil {
				fmt.Println("Option out of range!")
				continue
			}
			return str, m[str[0]], nil
		}
	}
}

// hfsf is a handler func setup func. See the default above
func newClient(l *logger.Logger, hfsf func(*ui.Handler)) *client {
	h := ui.NewHandler(os.Stdin, bufio.ScanLines)
	hfsf(h)
	return &client{
		Screen: ui.New(h),
		Logger: l,
	}
}
