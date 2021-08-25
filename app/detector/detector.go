package detector

import "github.com/7phs/area-51/app/lib"

var (
	_ Detector = (*detector)(nil)
)

type Detector interface {
	Start()
	Stop()
}

type detector struct {
	stream   RecordStream
	shutdown lib.Shutdown
}

func NewDetector(stream RecordStream) Detector {
	return &detector{
		stream:   stream,
		shutdown: lib.NewShutdown(),
	}
}

func (d *detector) Start() {
	d.shutdown.Add(1)
	go func() {
		defer d.shutdown.Done()

		for range d.stream.Records() {
		}
	}()
}

func (d *detector) Stop() {
	d.shutdown.Stop(nil, nil)
}
