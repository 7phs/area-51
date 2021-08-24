package watcher

import (
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
)

var (
	_ Watcher = (*watcher)(nil)
)

type Watcher interface {
	WatchFileChanges(filePath string) (Queue, error)

	Start()
	Stop()
}

type watcher struct {
	sync.RWMutex

	notifier *fsnotify.Watcher

	watchedDir map[string]bool
	queues     map[string]Queue

	wg       sync.WaitGroup
	shutdown chan bool
	once     sync.Once
}

func NewWatcher() (Watcher, error) {
	n, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, ErrUnexpected(err)
	}

	return &watcher{
		notifier:   n,
		watchedDir: make(map[string]bool),
		queues:     make(map[string]Queue),
		shutdown:   make(chan bool),
	}, nil
}

// WatchFileChanges - a method is responsible to create a queue to watch file changes
// A queue is identified by file name. File name of watched files should be unique
// TODO: needs implement new of identification and watching different directories
func (w *watcher) WatchFileChanges(filePath string) (Queue, error) {
	// check shutdown
	select {
	case <-w.shutdown:
		return nil, ErrAlreadyShutdown()
	default:
	}

	queue := newQueue(filePath)

	w.Lock()
	defer w.Unlock()

	if _, ok := w.watchedDir[queue.Dir()]; !ok {
		if err := w.notifier.Add(queue.Dir()); err != nil {
			return nil, ErrUnexpected(err)
		}
	}

	w.queues[queue.FileName()] = queue

	return queue, nil
}

func (w *watcher) Start() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()

		w.reactor()
	}()
}

func (w *watcher) Stop() {
	w.once.Do(func() {
		close(w.shutdown)
		w.wg.Wait()
	})
}

func (w *watcher) reactor() {
	for {
		select {
		case <-w.shutdown:
			return

		case event := <-w.notifier.Events:
			w.eventHandler(event)

		case err := <-w.notifier.Errors:
			log.Println("event: ", err)
		}
	}
}

func (w *watcher) eventHandler(event fsnotify.Event) {
	q, ok := w.getQueue(event.Name)
	if !ok {
		return
	}

	switch event.Op {
	case fsnotify.Create:
		q.Send(Open)

	case fsnotify.Write:
		q.Send(Read)

	case fsnotify.Remove,
		fsnotify.Rename:
		q.Send(Close)
	}
}

func (w *watcher) getQueue(fileName string) (Queue, bool) {
	w.RLock()
	defer w.RUnlock()

	q, ok := w.queues[fileName]
	return q, ok
}
