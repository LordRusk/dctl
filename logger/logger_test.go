package logger

import (
	"io"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	loggerTestFile = Home(loggerTestFile)
	l, err := New(io.Discard, "", 0, loggerTestFile)
	if err != nil {
		t.Errorf("Failed to create logger: %s\n", err)
	}
	defer l.Close()
	defer func() {
		if err := os.RemoveAll(Parent(loggerTestFile)); err != nil {
			t.Errorf("Failed to remove test log file: %s\n", err)
		}
	}()
	l.Print(loggerTestText)
	bites, err := os.ReadFile(loggerTestFile)
	if err != nil {
		t.Errorf("Failed to read test file: %s\n", err)
	}

	if string(bites) != loggerTestText {
		t.Errorf("Failed! Logged text and test output file text different!\nLogged: \"%s\"\nRead: \"%s\"\n", loggerTestText, bites)
	}
}

var loggerTestFile = "~/tmpTestDir/testLogFile.log"
var loggerTestText = "These are the test lines\nThese lines are testing lines\nuriel was right\nThese are some more test lines with 😀😀differencees\nGod is dead. God remains dead. And we have killed him. How shall we comfort ourselves, the murderers of all murderers? What was holiest and mightiest of all that the world has yet owned has bled to death under our knives: who will wipe this blood off us? What water is there for us to clean ourselves? What festivals of atonement, what sacred games shall we have to invent? Is not the greatness of this deed too great for us? Must we ourselves not become gods simply to appear worthy of it?\n\n"
