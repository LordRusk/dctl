package ui

import (
	"bufio"
	"os"
	"testing"
)

func TestHandlerBasic(t *testing.T) {
	h := NewHandler(os.Stdin, bufio.ScanLines)
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
}

// I wrote this for a bug that didn't exist...
func TestHandlerInput(t *testing.T) {
	var i int
	for k, v := range thInputMenu {
		if thInputMenu[k] != v {
			t.Error("How did we get here?")
		}
		if thInputMenu[thiInputs[i]] != v {
			t.Errorf("Keys are not the same: '%s' '%s'", thInputMenu[thiInputs[i]], k)
		}
		i++
	}
}

var thBasicNeededStr = "needed str"
var thBasicMenu = Menu{
	0: "none big man",
	1: "uno big woman",
	2: "kys dumb one",
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
var thInputMenu = Menu{
	0:              "num input!",
	"String Input": 432,
	"🌡🌡emojis 🔥🔥":  "Sucks for people for can't insert emoji's",
}
