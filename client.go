package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/diamondburned/arikawa/v2/state"
	"github.com/lordrusk/dctl/logger"
	"github.com/lordrusk/dctl/menu"
	"github.com/lordrusk/dctl/ui"
)

type client struct {
	*state.State
	*ui.Handler
	*logger.Logger
}

// returns the keys of a ui.Menu
// as a slice of strings
func keysToSlice(m ui.Menu) []string {
	strs := make([]string, len(m))
	var i int
	for k := range m {
		strs[i] = fmt.Sprintf("%v", k)
		i++
	}
	return strs
}

// returns the values of a ui.Menu
// as a slice of strings
func valuesToSlice(m ui.Menu) []string {
	strs := make([]string, len(m))
	var i int
	for _, v := range m {
		strs[i] = fmt.Sprintf("%v", v)
		i++
	}
	return strs
}

func setupDefaultHandlerFunc(h *ui.Handler) {
	// key: string, value: interface{}
	h.Hf = func(m ui.Menu) (interface{}, interface{}, error) { // ui.HandlerFunc
		// menu.Map building strung in
		var b strings.Builder
		mapp := make(menu.Map)
		for k, v := range m {
			str := fmt.Sprintf("[%s] %s\n", k, v)
			b.WriteString(str)
			mapp[str] = [2]interface{}{k, v}
		}

		for {
			if *useDmenu {
				d, err := menu.NewDmenu(b.String(), defaultDmenuOpts, "-P", mapp)
				if err != nil {
					fmt.Printf("Failed to run menu: %s\n", err)
				}

				if err := d.Prompt(false); err != nil {
					fmt.Println(err)
					continue
				}

				kav := d.Map()[string(d.Bytes())+"\n"]
				for k, v := range m {
					if kav[0] == k && kav[1] == v {
						return k, v, nil
					}
				}

				fmt.Println("Option out of range!")
			} else {
				str, err := h.NextString(b.String())
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
}

// hfsf is a handler func setup func. See the default above
func newClient(l *logger.Logger, hfsf func(*ui.Handler), passFlag string) *client {
	// setup client
	h := ui.NewHandler(os.Stdin, bufio.ScanLines, defaultDmenuOpts, passFlag, *useDmenu)
	hfsf(h)

	return &client{
		Handler: h,
		Logger:  l,
	}
}
