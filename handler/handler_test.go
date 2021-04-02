package handler

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// This is not a good test and it could fail
// because of hardware (potentially). If it
// ever does, I'll write a better one then.
func TestHandler(t *testing.T) {
	var strs []string
	h := New(func(killCh chan struct{}, _ chan error) {
		for {
			time.Sleep(time.Duration(50) * time.Millisecond)
			strs = append(strs, "Another str!")
			select {
			case <-killCh:
				return
			default:
			}
		}
	}, func(err error) {
		fmt.Println(err)
	})
	if err := h.Stop(); err == nil { // this should return an error
		t.Error("(*Handler).Stop() should return an error when handler is not running!")
	}
	h.Start()
	if err := h.Start(); err == nil { // this should return an error
		t.Error("(*Handler).Start() should return an error when handler is not running!")
	}
	time.Sleep(time.Duration(1) * time.Second)
	h.Stop()
	if len(strs) != 19 {
		t.Error("Number os strings should be 19....so make a better fucking test dipshit")
	}
	h.Start()
	time.Sleep(time.Duration(20) * time.Millisecond)
	h.Stop()
	if len(strs) != 20 {
		t.Error("Write a better fucking test")
	}
}

func TestHandlerErrors(t *testing.T) {
	metaErrCh := make(chan error)
	h := New(func(_ chan struct{}, errCh chan error) {
		_, err := os.Open(time.Now().Format(time.UnixDate)) // this file shouldn't exist
		errCh <- err
	}, func(err error) {
		metaErrCh <- err
	})
	h.Start()
	err := <-metaErrCh
	if err == nil {
		t.Error("returned error should not be nil")
	}
	if h.Running() {
		t.Error("handler should not be running after error!")
	}
}
