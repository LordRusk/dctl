package logger

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
)

var home string

func init() { home, _ = os.UserHomeDir() }

type Logger struct {
	*log.Logger
	Close func() error // closes the file
}

// replaces '~' with '/home/[name]'
// if home == "" then it does nothing
func Home(path string) string {
	if home != "" {
		return strings.ReplaceAll(path, "~", home)
	}
	return path
}

// returns the parent directories of a path
func Parent(path string) string {
	strs := strings.Split(path, "/")
	return strings.Join(strs[:len(strs)-1], "/")
}

// Allows for logging to an io.Writer and a file easily
func New(out io.Writer, prefix string, flag int, path string) (*Logger, error) {
	path = Home(path)
	if err := os.MkdirAll(Parent(path), 0777); err != nil {
		return nil, errors.Wrap(err, "Failed creating log file")
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, errors.Wrap(err, "could not create file")
	}
	mwr := io.MultiWriter(out, f)
	return &Logger{
		Logger: log.New(mwr, prefix, flag),
		Close: func() error {
			if err := f.Close(); err != nil {
				return errors.Wrap(err, "unable to close file")
			}
			// remove the log file if nothing was logged
			bites, err := os.ReadFile(path)
			if err != nil {
				return errors.Wrap(err, "unable to read file")
			}
			if len(bites) < 1 {
				if err := os.Remove(path); err != nil {
					return errors.Wrap(err, "unable to remove file")
				}
			}
			return nil
		},
	}, nil
}
