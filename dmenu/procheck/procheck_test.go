package procheck

import (
	"testing"
	"time"
)

func TestIsInstalled(t *testing.T) {
	if !IsInstalled("go") { // go should be installed
		t.Error("How is go not installed?")
	}
	now := time.Now().Format(time.UnixDate)
	if IsInstalled(now) { // a binarie of this name shouldn't exist
		t.Errorf("How does '%s' exist?", now)
	}
}
