package esaySpinner

import (
	"fmt"
	"io"
	"sync"
	"time"
	"unicode/utf8"
)

type Spinner struct {
	mu         *sync.RWMutex //
	chars      []string      // chars holds the chosen character set
	active     bool          // active holds the state of the spinner
	stopChan   chan struct{} // stopChan is a channel used to stop the indicator
	lastOutput string        // last character(set) written
	*Options
}

// Option is a function that takes a spinner and applies
// a given configuration.
type Option func(*Spinner)

// Options contains fields to configure the spinner.
type Options struct {
	Delay      time.Duration // Delay is the speed of the indicator
	FinalMSG   string
	Writer     io.Writer // to make testing better, exported so users have access. Use `WithWriter` to update after initialization.
	HideCursor bool      // hideCursor determines if the cursor is visible
}

// New provides a pointer to an instance of Spinner with the supplied options.
func New(cs []string, m *sync.RWMutex, d time.Duration, io io.Writer, options ...Option) *Spinner {
	s := &Spinner{
		chars:    cs,
		mu:       m,
		active:   false,
		stopChan: make(chan struct{}, 1),
	}
	s.Options = &Options{
		Delay:  d,
		Writer: io,
	}

	for _, option := range options {
		option(s)
	}
	return s
}

// WithHiddenCursor hides the cursor
// if hideCursor = true given.
func WithHiddenCursor(hideCursor bool) Option {
	return func(s *Spinner) {
		s.HideCursor = hideCursor
	}
}

func (s *Spinner) WithHiddenCursor(hideCursor bool) *Spinner {
	s.HideCursor = hideCursor
	return s
}

// WithWriter adds the given writer to the spinner. This
// function should be favored over directly assigning to
// the struct value.
func WithWriter(w io.Writer) Option {
	return func(s *Spinner) {
		s.mu.Lock()
		s.Writer = w
		s.mu.Unlock()
	}
}

func (s *Spinner) WithWriter(w io.Writer) *Spinner {
	s.Writer = w
	return s
}

// WithFinalMSG adds the given string ot the spinner
// as the final message to be written.
func WithFinalMSG(finalMsg string) Option {
	return func(s *Spinner) {
		s.FinalMSG = finalMsg
	}
}

func (s *Spinner) WithFinalMSG(finalMsg string) *Spinner {
	s.FinalMSG = finalMsg
	return s
}

// Active will return whether or not the spinner is currently active.
func (s *Spinner) Active() bool {
	return s.active
}

// Start will start the indicator.
func (s *Spinner) Start() {
	s.mu.Lock()
	if s.active {
		s.mu.Unlock()
		return
	}
	if s.HideCursor {
		// hides the cursor
		fmt.Fprint(s.Writer, "\033[?25l")
	}
	s.active = true
	s.mu.Unlock()

	go func() {
		for {
			for i := 0; i < len(s.chars); i++ {
				select {
				case <-s.stopChan:
					return
				default:
					s.mu.Lock()
					if !s.active {
						s.mu.Unlock()
						return
					}
					outPlain := fmt.Sprintf("%s", s.chars[i])
					s.lastOutput = s.chars[i]
					fmt.Fprint(s.Writer, outPlain)
					s.erase()
					delay := s.Delay
					s.mu.Unlock()
					time.Sleep(delay)
				}
			}
		}
	}()

}

// Stop stops the indicator.
func (s *Spinner) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.active {
		s.active = false
		if s.HideCursor {
			// makes the cursor visible
			fmt.Fprint(s.Writer, "\033[?25h")
		}
		s.erase()
		if s.FinalMSG != "" {
			fmt.Fprint(s.Writer, s.FinalMSG)
		}
		s.stopChan <- struct{}{}
	}
}

// erase deletes written characters.
// Caller must already hold s.lock.
func (s *Spinner) erase() {
	n := utf8.RuneCountInString(s.lastOutput)
	//if n == 0 {
	//	fmt.Fprintf(s.Writer, "\033[K") // erases to end of line from the position of cursor
	//} else {
	//	fmt.Fprintf(s.Writer, "\033[%dD\033[K", n) // erases to end of line from beginning of cursor
	//}
	if n != 0 {
		fmt.Fprintf(s.Writer, "\033[%dD", n)
	}
	s.lastOutput = ""
}

// Lock allows for manual control to lock the spinner.
func (s *Spinner) Lock() {
	s.mu.Lock()
}

// Unlock allows for manual control to unlock the spinner.
func (s *Spinner) Unlock() {
	s.mu.Unlock()
}
