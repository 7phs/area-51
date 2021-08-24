package watcher

const (
	Open Event = iota
	Read
	Close
)

type Event int
