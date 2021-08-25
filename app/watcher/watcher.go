package watcher

import (
	"log"
	"path/filepath"
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
func (w *watcher) WatchFileChanges(filePath string) (Queue, error) {
	// check shutdown
	select {
	case <-w.shutdown:
		return nil, ErrAlreadyShutdown()
	default:
	}

	q := newQueue(filePath)

	w.Lock()
	defer w.Unlock()

	dir := filepath.Dir(filePath)
	filePath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, ErrUnexpected(err)
	}

	if _, ok := w.watchedDir[dir]; !ok {
		if err := w.notifier.Add(dir); err != nil {
			return nil, ErrUnexpected(err)
		}
	}

	w.queues[filePath] = q

	return q, nil
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
			w.handleEvent(event)

		case err := <-w.notifier.Errors:
			log.Println("error on event handling: ", err)
		}
	}
}

func (w *watcher) handleEvent(event fsnotify.Event) {
	filePath, err := filepath.Abs(event.Name)
	if err != nil {
		return
	}

	q, ok := w.getQueue(filePath)
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
