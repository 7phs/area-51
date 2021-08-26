package detector

import (
	"log"
	"time"

	data_stream "github.com/7phs/area-51/app/data-stream"
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
			loaded     = false
			stat       PartitionStat
			ok         bool
		)

		for rec := range d.stream.Records() {
			if rec.IsCommand() {
				switch rec.Command {
				case data_stream.EOF, data_stream.CloseData:
					log.Println(time.Now(), "DETECTOR: ", totalCount, " / ", time.Since(start))

					d.cleanWriter.Flush()
					d.anomaliesWriter.Flush()

					totalCount = 0
					start = time.Now()
				}

				continue
			}

			// wait for stat
			if !loaded {
				stat, ok = <-d.validator.Validator()
				loaded = true
				start = time.Now()
			} else {
				select {
				case stat, ok = <-d.validator.Validator():
				default:
				}
			}

			if !ok {
				log.Println(time.Now(), "stat updated is closed")
				return
			}

			totalCount++

			d.writeRecord(rec, stat.IsAnomaly(string(rec.Key), rec.FeaturesF64))

			if totalCount >= 50_000 {
				log.Println(time.Now(), "DETECTOR: ", totalCount, " / ", time.Since(start))

				totalCount = 0
				start = time.Now()
			}
		}
	}()

	d.stream.Start()
}

func (d *detector) writeRecord(rec DataRecord, isAnomaly bool) {
	if !isAnomaly {
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
