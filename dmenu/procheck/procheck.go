// procheck is a package to check
// if a system has a given binary
package procheck

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// ---------------------name---exists
var binaries = make(map[string]bool)

func init() {
	path := strings.Split(os.Getenv("PATH"), ":")
	for _, dir := range path {
		bs, err := exec.Command("ls", dir).Output()
		if err != nil {
			panic(errors.Wrap(err, fmt.Sprintf("could not ls directory '%s'", dir)))
		}
		for _, binary := range strings.Split(string(bs), "\n") {
			binaries[strings.TrimSpace(binary)] = true
		}
	}
}

func IsInstalled(proname string) bool {
	return binaries[proname]
}
