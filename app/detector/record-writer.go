package detector

import (
	"bufio"
	"log"
	"os"
	"time"
)

const (
	defaultBufSize = 5 * 1024 * 1024
	header         = `id,feature,city,sport,size
`
)

type RecordWriter interface {
	Write(record DataRecord)
	Flush()
	Close()
}

type recordWriter struct {
	fileName string
	file     *os.File
	buf      *bufio.Writer
}

func NewRecordWriter(fileName string) RecordWriter {
	return &recordWriter{
		fileName: fileName,
		buf:      bufio.NewWriterSize(nil, defaultBufSize),
	}
}

func (w *recordWriter) init() {
	if w.file != nil {
		return
	}

	var err error

	log.Println(time.Now(), "open file '"+w.fileName+"' to write data")

	w.file, err = os.OpenFile(w.fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Println(time.Now(), "failed to open file '"+w.fileName+"' to write data: ", err)
		return
	}

	stat, err := w.file.Stat()
	if err != nil {
		log.Println(time.Now(), "failed to get stat of file '"+w.fileName+"': ", err)
		return
	}

	w.buf.Reset(w.file)

	if stat.Size() == 0 {
		if _, err := w.buf.Write([]byte(header)); err != nil {
			log.Println(time.Now(), "failed to write header to file '"+w.fileName+"': ", err)
		}
	}
}

func (w *recordWriter) Write(record DataRecord) {
	w.init()

	record.Serialize(w.buf)
}

func (w *recordWriter) Flush() {
	if w.file == nil {
		return
	}

	log.Println(time.Now(), "flush data to file '"+w.fileName+"'")

	if err := w.buf.Flush(); err != nil {
		log.Println(time.Now(), "failed to flush data to file '"+w.fileName+"': ", err)
	}
}

func (w *recordWriter) Close() {
	if w.file == nil {
		return
	}

	log.Println(time.Now(), "close file '"+w.fileName+"'")

	if err := w.buf.Flush(); err != nil {
		log.Println(time.Now(), "failed to flush data to file '"+w.fileName+"': ", err)
	}

	if err := w.file.Close(); err != nil {
		log.Println(time.Now(), "failed to close file '"+w.fileName+"': ", err)
	}

	w.buf.Reset(nil)

	w.file = nil
}
