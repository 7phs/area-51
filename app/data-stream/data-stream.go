package data_stream

import (
	"io"

	"github.com/7phs/area-51/app/watcher"
)

var (
	_ DataStream = (*dataStream)(nil)
)

type DataStream interface {
	io.ReadCloser

	Start()
	Stop()
}

type dataStream struct {
}

func NewDataStream(queue watcher.Queue) (DataStream, error) {

}
