package detector

import (
	"github.com/7phs/area-51/app/lib"
	"log"
	"time"
)

var (
	_ Detector = (*detector)(nil)
)

type Detector interface {
	Start()
	Stop()
}

type detector struct {
	stream   RecordReader
	shutdown lib.Shutdown
}

func NewDetector(stream RecordReader) Detector {
	return &detector{
		stream:   stream,
		shutdown: lib.NewShutdown(),
	}
}

func (d *detector) Start() {
	d.shutdown.Add(1)
	go func() {
		defer d.shutdown.Done()

		totalCount := int64(0)
		start := time.Now()

		for range d.stream.Records() {
			totalCount++

			if totalCount > 10_000 {
				log.Println("10 000 per ", time.Since(start))

				totalCount = 0
				start = time.Now()
			}
		}
	}()

	d.stream.Start()
}

func (d *detector) Stop() {
	d.stream.Stop()

	d.shutdown.Stop(nil, nil)
}
