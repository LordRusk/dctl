// this is a "simple" abstraction making it easy to
// handle arbitrary functions concurrently.
//
// i really wish i could make this simpler.
package handler

import "github.com/pkg/errors"

type HandlerFunc func(chan struct{}, chan error)
type ErrorHandlerFunc func(error)

type Handler struct {
	hf      HandlerFunc
	ehf     ErrorHandlerFunc
	killErr chan struct{} // kills the error handler loop
	killHan chan struct{} // kills the handler loop
	errCh   chan error
	running bool
}

// if a concurrent handler errors, it is expected to
// end, using (*handler).Start() to start it again.
// problems will arrise if you do not ensure this.
func New(hf HandlerFunc, ehf ErrorHandlerFunc) *Handler {
	killErr := make(chan struct{})
	killHan := make(chan struct{})
	errCh := make(chan error)

	return &Handler{
		hf:      hf,
		ehf:     ehf,
		killErr: killErr,
		killHan: killHan,
		errCh:   errCh,
	}
}

// this will only return an error if
// the handler is already running.
func (h *Handler) Start() error {
	if h.running {
		return errors.New("Handler already running!")
	}
	go func() {
		select {
		case err := <-h.errCh:
			h.ehf(err)
			h.running = false
		case <-h.killErr:
			h.killHan <- struct{}{}
			h.running = false
			return
		}
	}()
	go h.hf(h.killHan, h.errCh)
	h.running = true
	return nil
}

// this will only return ann error if
// the handler isn't running.
func (h *Handler) Stop() error {
	if !h.running {
		return errors.New("Handler already running!")
	}
	h.killErr <- struct{}{}
	h.running = false
	return nil
}

func (h *Handler) Running() bool {
	return h.running
}
