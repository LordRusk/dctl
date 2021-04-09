// dmenu is a package for simple dmenu integrations
// in go programs.
package dmenu

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	"github.com/lordrusk/dctl/dmenu/procheck"
)

// pkgReady is whether this package
// can be used on the current system
var pkgReady = procheck.IsInstalled("dmenu")

// str is the opts seperated by '\n'
// flags are the array of flags passed to dmenu
func Prompt(str string, flags []string) (string, error) {
	if !pkgReady {
		return "", errors.New("dmenu not installed!")
	}
	cmd := exec.Command("dmenu", flags[:]...)
	cmd.Stdin = strings.NewReader(str)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("unable to run dmenu: %s", output))
	}
	return strings.TrimSpace(string(output)), nil
}

// Prompt, but takes a []string
func PromptSlice(strs []string, flags []string) (string, error) {
	return Prompt(strings.Join(strs, "\n"), flags)
}
