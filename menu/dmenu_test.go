// dmenu must be installed
package menu

import (
	"flag"
	"fmt"
	"os/exec"
	"testing"
)

// dmenu needs a specific patch
// for this test not to fail
var doPassTest = flag.Bool("P", false, "Do Password Tests")
var passOpt = "-P"

func TestDmenu(t *testing.T) {
	flag.Parse()

	if _, err := exec.LookPath("dmenu"); err != nil {
		fmt.Println("Dmenu not installed...")
		return
	}

	d, err := NewSlicedDmenu([]string{"choice one", "choice two", "choice three"}, []string{"-l", "10"}, "", nil)
	if err != nil {
		t.Fatal(err)
	}

	if d.Ran() {
		t.Fatal("d.Ran should not be true")
	}
	if err := d.Prompt(false); err != nil {
		fmt.Println(err)
	}
	if !d.Ran() {
		t.Fatal("d.Ran should not be false")
	}
	if err := d.Prompt(false); err == nil {
		t.Fatal("Menu should not be runable again")
	}
	if d.Bytes() == nil {
		t.Error("d.Bytes() should only be nil if user escapes")
	}
	fmt.Printf("%s", d.Bytes())

	// password test
	if !*doPassTest {
		return
	}

	d, err = NewSlicedDmenu([]string{"choice one", "choice two", "choice three"}, []string{"-l", "10"}, passOpt, nil)
	if err != nil {
		t.Fatal(err)
	}

	if err := d.Prompt(true); err != nil {
		// t.Error("No selection or password functionality not patched into dmenu")
		t.Error(err)
	}

	if d.Bytes() == nil {
		fmt.Println("Password bytes are nil!")
		return
	}
	fmt.Printf("%s", d.Bytes())
}
