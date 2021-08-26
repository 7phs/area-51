package data_stream

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"sync"
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
	Read() <-chan Buffer
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
	reader   chan Buffer
	cmd      chan Command

	readWG   sync.WaitGroup
	shutdown lib.Shutdown
}

func NewDataStream(queue FileChangesQueue) DataStream {
	return &dataStream{
		queue:    queue,
		filePath: queue.FilePath(),
		f:        nil,
		buf:      bufio.NewReaderSize(nil, defaultBufSize),
		reader:   make(chan Buffer),
		cmd:      make(chan Command),
		shutdown: lib.NewShutdown(),
	}
}

func (d *dataStream) Read() <-chan Buffer {
	return d.reader
}

func (d *dataStream) Start() {
	d.shutdown.Add(1)
	go func() {
		defer d.shutdown.Done()

		d.reactor()
	}()

	d.readWG.Add(1)
	go func() {
		defer d.readWG.Done()

		d.readStream()
	}()
}

func (d *dataStream) Stop() {
	d.shutdown.Stop(nil, func() {
		close(d.cmd)

		d.readWG.Wait()

		close(d.reader)
	})

	if d.f != nil {
		d.f.Close()
		d.f = nil
	}
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

		d.sendCmd(NewData)

	case Read:
		if d.f == nil {
			return
		}

		d.sendCmd(Data)

	case Close:
		if d.f == nil {
			return
		}

		d.reader <- CloseBuffer(nil)

		log.Println(time.Now(), "close file: ", d.filePath)

		if err := d.f.Close(); err != nil {
			log.Println(time.Now(), "failed to close file: ", d.filePath)
		}

		d.f = nil
		d.buf.Reset(nil)
	}
}

func (d *dataStream) sendCmd(command Command) {
	go func() {
		d.cmd <- command
	}()
}

func (d *dataStream) readStream() {
	for cmd := range d.cmd {
		d.tryReadBuf(cmd)
	}
}

func (d *dataStream) tryReadBuf(cmd Command) {
	for {
		select {
		case <-d.shutdown.Ch():
			return
		default:
		}

		buf := make([]byte, defaultBufSize)

		n, err := d.buf.Read(buf)
		switch {
		case errors.Is(err, io.EOF):
			cmd = EOF

		case err != nil:
			log.Println(time.Now(), "error while read: ", err)
			return
		}

		buf = buf[:n]

		d.reader <- NewBuffer(cmd, buf)

		switch cmd {
		case NewData:
			cmd = Data
		case EOF:
			return
		}
	}
}
