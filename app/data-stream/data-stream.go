package data_stream

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"sync"

	"github.com/7phs/area-51/app/watcher"
)

const (
	defaultBufSize = 1024 * 1024
)

var (
	_ DataStream = (*dataStream)(nil)
)

type DataStream interface {
	Read() <-chan []byte
	Start()
	Stop()
}

type dataStream struct {
	queue    watcher.Queue
	filePath string
	f        *os.File
	buf      *bufio.Reader
	reader   chan []byte
	cmd      chan bool

	wg       sync.WaitGroup
	shutdown chan bool
	once     sync.Once
}

func NewDataStream(queue watcher.Queue) DataStream {
	return &dataStream{
		queue:    queue,
		filePath: queue.FilePath(),
		f:        nil,
		buf:      bufio.NewReaderSize(nil, defaultBufSize),
		reader:   make(chan []byte),
		cmd:      make(chan bool),
		shutdown: make(chan bool),
	}
}

func (d *dataStream) Read() <-chan []byte {
	return d.reader
}

func (d *dataStream) Start() {
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()

		d.reactor()
	}()

	go func() {
		d.readStream()
	}()
}

func (d *dataStream) Stop() {
	d.once.Do(func() {
		close(d.shutdown)
		d.wg.Wait()
		close(d.cmd)
		close(d.reader)
	})
}

func (d *dataStream) reactor() {
	for {
		select {
		case <-d.shutdown:
			return

		case event := <-d.queue.Ch():
			d.handleEvent(event)
		}
	}
}

func (d *dataStream) handleEvent(event watcher.Event) {
	switch event {
	case watcher.Open:
		var err error

		log.Println("open file: ", d.filePath)

		d.f, err = os.OpenFile(d.filePath, os.O_RDONLY, 0)
		if err != nil {
			log.Println("failed to open file '"+d.filePath+"': ", err)
			return
		}
		d.buf.Reset(d.f)

		d.readCmd()

	case watcher.Read:
		if d.f == nil {
			return
		}

		d.readCmd()

	case watcher.Close:
		if d.f == nil {
			return
		}

		log.Println("close file: ", d.filePath)

		if err := d.f.Close(); err != nil {
			log.Println("failed to close file: ", d.filePath)
		}

		d.f = nil
		d.buf.Reset(nil)
	}
}

func (d *dataStream) readCmd() {
	go func() {
		d.cmd <- true
	}()
}

func (d *dataStream) readStream() {
	for range d.cmd {
		d.tryReadBuf()
	}
}

func (d *dataStream) tryReadBuf() {
	for {
		buf := make([]byte, defaultBufSize)

		n, err := d.buf.Read(buf)
		if err != nil && !errors.Is(err, io.EOF) {
			log.Println("error while read: ", err)
		}
		if n == 0 {
			return
		}

		buf = buf[:n]

		go func() {
			d.reader <- buf
		}()
	}
}
