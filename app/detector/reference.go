package detector

import (
	"log"
	"time"

	"github.com/7phs/area-51/app/lib"
)

var (
	_ AnomaliesValidator = (*reference)(nil)
	_ Reference          = (*reference)(nil)
)

type AnomaliesValidator interface {
	Validate(rec DataRecord) bool
}

type Reference interface {
	AnomaliesValidator

	Start()
	Stop()
}

type reference struct {
	stream   RecordReader
	stat     PartitionStat
	shutdown lib.Shutdown
}

func NewReference(stream RecordReader) Reference {
	return &reference{
		stream:   stream,
		stat:     NewPartitionStat(),
		shutdown: lib.NewShutdown(),
	}
}

func (r *reference) Validate(_ DataRecord) bool {
	return true
}

func (r *reference) Start() {
	r.shutdown.Add(1)
	go func() {
		defer r.shutdown.Done()

		var (
			totalCount = int64(0)
			start      = time.Now()
		)

		for rec := range r.stream.Records() {
			totalCount++

			r.stat.Add(string(rec.Key), rec.FeaturesF64)

			if totalCount > 45_000 {
				log.Println("REFERENCE: 45 000 per ", time.Since(start))

				totalCount = 0
				start = time.Now()
			}
		}
	}()

	r.stream.Start()
}

func (d *reference) Stop() {
	d.stream.Stop()

	d.shutdown.Stop(nil, nil)
}
