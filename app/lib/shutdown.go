package lib

import (
	"sync"
)

type Shutdown struct {
	sync.WaitGroup

	done chan bool
	once sync.Once
}

func NewShutdown() Shutdown {
	return Shutdown{
		done: make(chan bool),
	}
}

func (s *Shutdown) Ch() <-chan bool {
	return s.done
}

func (s *Shutdown) Stop(pre, post func()) {
	s.once.Do(func() {
		if pre != nil {
			pre()
		}

		close(s.done)
		s.Wait()

		if post != nil {
			post()
		}
	})
}
