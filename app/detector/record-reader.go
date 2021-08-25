package detector

import (
	data_stream "github.com/7phs/area-51/app/data-stream"
	"github.com/7phs/area-51/app/lib"
)

var (
	_ RecordStream = (*recordReader)(nil)
	_ RecordReader = (*recordReader)(nil)
)

type RecordStream interface {
	Records() <-chan DataRecord
}

type RecordReader interface {
	RecordStream

	Start()
	Stop()
}

type recordReader struct {
	reader     data_stream.DataReader
	delimiter  byte
	skipHeader bool

	records chan DataRecord

	shutdown lib.Shutdown
}

func NewRecordReader(reader data_stream.DataReader) RecordReader {
	return &recordReader{
		reader:     reader,
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
}

func (p *recordReader) Stop() {
	p.shutdown.Stop(nil, nil)
}

func (p *recordReader) processor() {
	totalCount := 0
	for buf := range p.reader.Read() {
		prev := 0
		count := 0

		for i, c := range buf {
			if c != '\n' {
				// TODO: check maximum length of line
				continue
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
	go func() {
		p.records <- rec
	}()
}
