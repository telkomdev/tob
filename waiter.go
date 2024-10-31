package tob

// Waiter the waiter that follow the semaphore pattern
type Waiter interface {
	Done()
	Wait()
	Close()
}

type empty struct{}

type waiter struct {
	Capacity uint
	Sig      chan empty
}

// NewWaiter will return new Waiter
func NewWaiter(c uint) Waiter {
	sig := make(chan empty, c)
	return &waiter{Capacity: c, Sig: sig}
}

// Done will tell the waiter if the execution is done
func (w *waiter) Done() {
	e := empty{}
	w.Sig <- e
}

// Wait will wait all available execution
func (w *waiter) Wait() {
	var i uint
	for i = 0; i < w.Capacity; i++ {
		<-w.Sig
	}
}

// Close will close the Sig channel
func (w *waiter) Close() { close(w.Sig) }
