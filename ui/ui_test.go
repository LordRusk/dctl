package ui

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"testing"
)

var testMenu = flag.Bool("-M", false, "Use menu for tests")
var doPassTest = flag.Bool("-P", false, "Do password test")
var PassOpt = flag.String("-O", "-P", "Set the default password opt")

func TestHandlerBasic(t *testing.T) {
	flag.Parse()
	h := NewHandler(os.Stdin, bufio.ScanLines, nil, *PassOpt, *testMenu)
	h.Hf = func(m Menu) (interface{}, interface{}, error) {
		for key, value := range m {
			if value.(string) == thBasicNeededStr {
				return key, value, nil
			}
		}
		t.Error("Failed to find needed string!")
		return nil, nil, nil // <-- so it builds
	}

	if _, _, err := h.Hf(thBasicMenu); err != nil {
		t.Errorf("Failed handling menu: %s\n", err)
	}

	for i := 0; i < 2; i++ {
		if _, err := h.NextInt("give an int"); err != nil {
			t.Error(err)
		}
		if _, err := h.NextBytes("type something"); err != nil {
			t.Error(err)
		}
		if _, err := h.NextString("type something"); err != nil {
			t.Error(err)
		}

		if *doPassTest {
			str, err := h.PassInput("Enter Not A Password")
			if err != nil {
				t.Error(err)
			}
			fmt.Println(str)
		}
		h.UseMenu = !h.UseMenu
	}
}

var thBasicNeededStr = "needed str"
var thBasicMenu = Menu{
	0: "none big man",
	1: "uno big woman",
	2: "kys dumb two",
	3: thBasicNeededStr,
}

// this was done for a reason...
var thiFirstInput int = 0
var thiSecondInput string = "String Input"
var thiThirdInput string = "🌡🌡emojis 🔥🔥"
var thiInputs = []interface{}{
	thiFirstInput,
	thiSecondInput,
	thiThirdInput,
}
