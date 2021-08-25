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
	totalCount := 0
	firstLine := true

	for buf := range p.stream.Read() {
		if buf == nil {
			firstLine = true
			// TODO: handle finish of buffer
			continue
		}

		prev := 0
		count := 0

		for i, c := range buf {
			if c != '\n' {
				// TODO: check maximum length of line
				continue
			}

			if firstLine {
				firstLine = false

				if p.skipHeader {
					continue
				}
			}

			p.send(parseDataRecord(p.delimiter, buf[prev:i]))
			count++

			prev = i + 1
		}

		p.send(parseDataRecord(p.delimiter, buf[prev:]))

		count++
		totalCount += count
	}
}

func (p *recordReader) send(rec DataRecord) {
	p.shutdown.Add(1)
	go func() {
		defer p.shutdown.Done()

		p.records <- rec
	}()
}
