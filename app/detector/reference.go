package detector

import (
	"log"
	"time"

	data_stream "github.com/7phs/area-51/app/data-stream"
	"github.com/7phs/area-51/app/lib"
)

var (
	_ AnomaliesValidator = (*reference)(nil)
	_ Reference          = (*reference)(nil)
)

type AnomaliesValidator interface {
	Validator() <-chan PartitionStat
}

type Reference interface {
	AnomaliesValidator

	Start()
	Stop()
}

type reference struct {
	stream    RecordReader
	validator chan PartitionStat
	shutdown  lib.Shutdown
}

func NewReference(stream RecordReader) Reference {
	return &reference{
		stream:    stream,
		validator: make(chan PartitionStat),
		shutdown:  lib.NewShutdown(),
	}
}

func (r *reference) Validator() <-chan PartitionStat {
	return r.validator
}

func (r *reference) Start() {
	r.shutdown.Add(1)
	go func() {
		defer r.shutdown.Done()

		stat := NewPartitionStat()
		start := time.Now()
		count := 0

		for rec := range r.stream.Records() {
			if rec.IsCommand() {
				switch rec.Command {
				case data_stream.EOF, data_stream.CloseData:
					log.Println(time.Now(), "REFERENCE: ", count, " / ", time.Since(start))

					select {
					case r.validator <- stat:
					case <-r.shutdown.Ch():
						return
					}

					stat = NewPartitionStat()
					start = time.Now()
					count = 0
				}

				continue
			}

			count++
			stat.Add(string(rec.Key), rec.FeaturesF64)
		}
	}()

	r.stream.Start()
}

func (r *reference) Stop() {
	r.stream.Stop()

	r.shutdown.Stop(nil, func() {
		close(r.validator)
	})
}
