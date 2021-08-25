package data_stream

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"time"

	"github.com/7phs/area-51/app/lib"
)

const (
	defaultBufSize = 5 * 1024 * 1024
)

var (
	_ DataStream = (*dataStream)(nil)
)

type DataReader interface {
	Read() <-chan []byte
}

type DataStream interface {
	DataReader

	Start()
	Stop()
}

type dataStream struct {
	queue    FileChangesQueue
	filePath string
	f        *os.File
	buf      *bufio.Reader
	reader   chan []byte
	cmd      chan bool

	shutdown lib.Shutdown
}

func NewDataStream(queue FileChangesQueue) DataStream {
	return &dataStream{
		queue:    queue,
		filePath: queue.FilePath(),
		f:        nil,
		buf:      bufio.NewReaderSize(nil, defaultBufSize),
		reader:   make(chan []byte),
		cmd:      make(chan bool),
		shutdown: lib.NewShutdown(),
	}
}

func (d *dataStream) Read() <-chan []byte {
	return d.reader
}

func (d *dataStream) Start() {
	d.shutdown.Add(1)
	go func() {
		defer d.shutdown.Done()

		d.reactor()
	}()

	go func() {
		d.readStream()
	}()
}

func (d *dataStream) Stop() {
	d.shutdown.Stop(nil, func() {
		close(d.cmd)
		close(d.reader)
	})
}

func (d *dataStream) reactor() {
	for {
		select {
		case <-d.shutdown.Ch():
			return

		case event := <-d.queue.Ch():
			d.handleEvent(event)
		}
	}
}

func (d *dataStream) handleEvent(event Event) {
	switch event {
	case Open:
		var err error

		log.Println(time.Now(), "open file: ", d.filePath)

		d.f, err = os.OpenFile(d.filePath, os.O_RDONLY, 0)
		if err != nil {
			log.Println(time.Now(), "failed to open file '"+d.filePath+"': ", err)
			return
		}
		d.buf.Reset(d.f)

		d.readCmd()

	case Read:
		if d.f == nil {
			return
		}

		d.readCmd()

	case Close:
		if d.f == nil {
			return
		}

		log.Println(time.Now(), "close file: ", d.filePath)

		if err := d.f.Close(); err != nil {
			log.Println(time.Now(), "failed to close file: ", d.filePath)
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
			log.Println(time.Now(), "error while read: ", err)
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
