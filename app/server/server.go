package server

import (
	"log"

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

	referenceStream data_stream.DataStream
	referenceReader detector.RecordReader
	reference       detector.Reference

	rawStream data_stream.DataStream
	rawReader detector.RecordReader
	detector  detector.Detector

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

	// TODO: optimize definition
	referenceStream := data_stream.NewDataStream(referenceQueue)
	referenceReader := detector.NewRecordReader(referenceStream)
	reference := detector.NewReference(referenceReader)

	rawStream := data_stream.NewDataStream(rawQueue)
	rawReader := detector.NewRecordReader(rawStream)
	dtctr := detector.NewDetector(rawReader)

	return &server{
		w: w,

		referenceStream: referenceStream,
		referenceReader: referenceReader,
		reference:       reference,

		rawStream: rawStream,
		rawReader: rawReader,
		detector:  dtctr,

		shutdown: lib.NewShutdown(),
	}, nil
}

func (s *server) Start() {
	s.reference.Start()
	s.referenceReader.Start()
	s.referenceStream.Start()

	s.detector.Start()
	s.rawReader.Start()
	s.rawStream.Start()

	s.w.Start()

	s.Process()
}

func (s *server) Stop() {
	s.shutdown.Stop(func() {
		s.w.Stop()

		s.rawStream.Stop()
		s.rawReader.Stop()
		s.detector.Stop()

		s.referenceStream.Stop()
		s.referenceReader.Stop()
		s.reference.Stop()
	}, nil)
}

func (s *server) Process() {
	s.shutdown.Add(1)
	go func() {
		defer s.shutdown.Done()

		for {
			select {
			case <-s.shutdown.Ch():
				return

			case buf, ok := <-s.referenceStream.Read():
				if !ok {
					continue
				}

				log.Println("READ - REF: ", string(buf))
			}
		}
	}()
}
