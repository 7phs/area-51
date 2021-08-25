package watcher

var (
	_ Queue = (*queue)(nil)
)

type Queue interface {
	FilePath() string
	Send(Event)
	Ch() <-chan Event
}

type queue struct {
	path string
	ch   chan Event
}

func newQueue(filePath string) Queue {
	return &queue{
		path: filePath,
		ch:   make(chan Event),
	}
}

func (q *queue) FilePath() string {
	return q.path
}

func (q *queue) Send(e Event) {
	go func() {
		q.ch <- e
	}()
}

func (q *queue) Ch() <-chan Event {
	return q.ch
}
