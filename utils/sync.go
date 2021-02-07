/*
Copyright © 2021 Matthew Champagne <mmchamp95@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

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
}

// Wait waits for all processes to call Done and sends results to the provided callback.
// Waiting can be ended prematurely by returning an error in the callback.
func (e *ErrChan) Wait(cb func(err error) error) error {
	for i := 0; i < e.n; i++ {
		if err := cb(<-e.ch); err != nil {
			return err
		}
	}
	return nil
}
