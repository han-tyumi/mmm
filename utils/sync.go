package utils

// ErrChan supports waiting for multiple error returning processes.
type ErrChan struct {
	n  int
	ch chan error
}

// NewErrCh creates a new ErrCh to handle n error returning processes.
func NewErrCh(n int) *ErrChan {
	return &ErrChan{
		n:  n,
		ch: make(chan error),
	}
}

// Do signals the error channel with the return value of the passed function.
func (e *ErrChan) Do(fn func() error) {
	e.Done(fn())
}

// Done signals that one of the processes is done with the given error value.
func (e *ErrChan) Done(v error) {
	e.ch <- v
	e.n--
}

// Wait waits for all processes to call Done and sends results to the provided callback.
// Waiting can be ended prematurely by returning an error in the callback.
func (e *ErrChan) Wait(cb func(err error) error) error {
	for e.n > 0 {
		if err := cb(<-e.ch); err != nil {
			return err
		}
	}
	return nil
}
