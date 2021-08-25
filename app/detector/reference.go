package detector

import "github.com/7phs/area-51/app/lib"

var (
	_ Reference = (*reference)(nil)
)

type Reference interface {
	Start()
	Stop()
}

type reference struct {
	stream   RecordReader
	shutdown lib.Shutdown
}

func NewReference(stream RecordReader) Reference {
	return &reference{
		stream:   stream,
		shutdown: lib.NewShutdown(),
	}
}

func (d *reference) Start() {
	d.shutdown.Add(1)
	go func() {
		defer d.shutdown.Done()

		for range d.stream.Records() {
		}
	}()

	d.stream.Start()
}

func (d *reference) Stop() {
	d.stream.Stop()

	d.shutdown.Stop(nil, nil)
}
