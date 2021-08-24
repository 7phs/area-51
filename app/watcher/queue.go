package watcher

import "path/filepath"

var (
	_ Queue = (*queue)(nil)
)

type Queue interface {
	Dir() string
	FileName() string
	Send(Event)
	Ch() chan Event
}

type queue struct {
	dir      string
	fileName string
	ch       chan Event
}

func newQueue(filePath string) Queue {
	return &queue{
		dir:      filepath.Dir(filePath),
		fileName: filepath.Base(filePath),
		ch:       make(chan Event),
	}
}

func (q *queue) Dir() string {
	return q.dir
}

func (q *queue) FileName() string {
	return q.fileName
}

func (q *queue) Send(e Event) {
	go func() {
		q.ch <- e
	}()
}

func (q *queue) Ch() chan Event {
	return q.ch
}
