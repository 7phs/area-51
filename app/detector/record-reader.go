package detector

import (
	data_stream "github.com/7phs/area-51/app/data-stream"
	"github.com/7phs/area-51/app/lib"
)

var (
	_ RecordReader = (*recordReader)(nil)
)

type RecordReader interface {
	Records() <-chan DataRecord

	Start()
	Stop()
}

type recordReader struct {
	stream     data_stream.DataStream
	delimiter  byte
	skipHeader bool

	records chan DataRecord

	shutdown lib.Shutdown
}

func NewRecordReader(reader data_stream.DataStream) RecordReader {
	return &recordReader{
		stream:     reader,
		delimiter:  ',',
		skipHeader: true,
		records:    make(chan DataRecord),
		shutdown:   lib.NewShutdown(),
	}
}

func (p *recordReader) Records() <-chan DataRecord {
	return p.records
}

func (p *recordReader) Start() {
	p.shutdown.Add(1)
	go func() {
		defer p.shutdown.Done()

		p.processor()
	}()

	p.stream.Start()
}

func (p *recordReader) Stop() {
	p.stream.Stop()

	p.shutdown.Stop(nil, func() {
		close(p.records)
	})
}

func (p *recordReader) processor() {
	var (
		prevBuf   data_stream.Buffer
		firstLine = true
		prevIndex = -1
	)

	for buf := range p.stream.Read() {
		switch buf.Command() {
		case data_stream.NewData, data_stream.CloseData:
			firstLine = true
			prevBuf = nil
			prevIndex = -1
		}

		firstLine, prevIndex = parseCSV(p.delimiter, p.skipHeader, firstLine, prevIndex, prevBuf, buf, p.send)
		prevBuf = buf
	}
}

func (p *recordReader) send(rec DataRecord) {
	p.shutdown.Add(1)
	go func() {
		defer p.shutdown.Done()

		p.records <- rec
	}()
}
