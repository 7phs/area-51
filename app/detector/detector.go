package detector

import (
	"log"
	"time"

	"github.com/7phs/area-51/app/lib"
)

var (
	_ Detector = (*detector)(nil)
)

type Detector interface {
	Start()
	Stop()
}

type detector struct {
	stream    RecordReader
	validator AnomaliesValidator

	cleanWriter     RecordWriter
	anomaliesWriter RecordWriter

	shutdown lib.Shutdown
}

func NewDetector(
	stream RecordReader,
	validator AnomaliesValidator,
	cleanWriter RecordWriter,
	anomaliesWriter RecordWriter,
) Detector {
	return &detector{
		stream:    stream,
		validator: validator,

		cleanWriter:     cleanWriter,
		anomaliesWriter: anomaliesWriter,

		shutdown: lib.NewShutdown(),
	}
}

func (d *detector) Start() {
	d.shutdown.Add(1)
	go func() {
		defer d.shutdown.Done()

		var (
			totalCount = int64(0)
			start      = time.Now()
		)

		for rec := range d.stream.Records() {
			totalCount++

			d.writeRecord(rec)

			if totalCount > 45_000 {
				log.Println("DETECTOR: 45 000 per ", time.Since(start))

				totalCount = 0
				start = time.Now()
			}
		}
	}()

	d.stream.Start()
}

func (d *detector) writeRecord(rec DataRecord) {
	if d.validator.Validate(rec) {
		d.cleanWriter.Write(rec)
		return
	}

	d.anomaliesWriter.Write(rec)
}

func (d *detector) Stop() {
	d.stream.Stop()

	d.shutdown.Stop(nil, nil)

	d.cleanWriter.Close()
	d.anomaliesWriter.Close()
}
