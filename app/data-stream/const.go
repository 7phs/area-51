package data_stream

const (
	Open Event = iota
	Read
	Close
)

type Event int
