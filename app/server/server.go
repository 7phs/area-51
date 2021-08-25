package server

import (
	"github.com/7phs/area-51/app/config"
	data_stream "github.com/7phs/area-51/app/data-stream"
	"github.com/7phs/area-51/app/detector"
	"github.com/7phs/area-51/app/lib"
)

type Server interface {
	Start()
	Stop()
}

type server struct {
	w data_stream.Watcher

	reference detector.Reference

	detector detector.Detector

	shutdown lib.Shutdown
}

func New(conf config.Config) (Server, error) {
	w, err := data_stream.NewWatcher()
	if err != nil {
		return nil, ErrUnexpected("failed to initialize data-stream", err)
	}

	referenceQueue, err := w.WatchFileChanges(conf.ReferenceFile)
	if err != nil {
		return nil, ErrUnexpected("failed to initialize queue of reference", err)
	}

	rawQueue, err := w.WatchFileChanges(conf.RawDataFile)
	if err != nil {
		return nil, ErrUnexpected("failed to initialize queue of raw data file", err)
	}

	referenceStream := data_stream.NewDataStream(referenceQueue)
	referenceReader := detector.NewRecordReader(referenceStream)
	reference := detector.NewReference(referenceReader)

	rawStream := data_stream.NewDataStream(rawQueue)
	rawReader := detector.NewRecordReader(rawStream)
	dtctr := detector.NewDetector(rawReader)

	return &server{
		w: w,

		reference: reference,
		detector:  dtctr,

		shutdown: lib.NewShutdown(),
	}, nil
}

func (s *server) Start() {
	s.reference.Start()

	s.detector.Start()

	s.w.Start()
}

func (s *server) Stop() {
	s.shutdown.Stop(func() {
		s.w.Stop()

		s.detector.Stop()

		s.reference.Stop()
	}, nil)
}
