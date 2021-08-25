package server

import (
	"log"

	"github.com/7phs/area-51/app/config"
	data_stream "github.com/7phs/area-51/app/data-stream"
	"github.com/7phs/area-51/app/lib"
	"github.com/7phs/area-51/app/watcher"
)

type Server interface {
	Start()
	Process()
	Stop()
}

type server struct {
	w               watcher.Watcher
	referenceStream data_stream.DataStream
	rawStream       data_stream.DataStream

	shutdown lib.Shutdown
}

func New(conf config.Config) (Server, error) {
	w, err := watcher.NewWatcher()
	if err != nil {
		return nil, err
	}

	referenceQueue, err := w.WatchFileChanges(conf.ReferenceFile)
	if err != nil {
		return nil, err
	}

	rawQueue, err := w.WatchFileChanges(conf.RawDataFile)
	if err != nil {
		return nil, err
	}

	return &server{
		w:               w,
		referenceStream: data_stream.NewDataStream(referenceQueue),
		rawStream:       data_stream.NewDataStream(rawQueue),
		shutdown:        lib.NewShutdown(),
	}, nil
}

func (s *server) Start() {
	s.referenceStream.Start()
	s.rawStream.Start()
	s.w.Start()

	s.Process()
}

func (s *server) Stop() {
	s.shutdown.Stop(func() {
		s.w.Stop()
		s.referenceStream.Stop()
		s.rawStream.Stop()
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

			case buf, ok := <-s.rawStream.Read():
				if !ok {
					continue
				}

				log.Println("READ - RAW: ", string(buf))

			case buf, ok := <-s.referenceStream.Read():
				if !ok {
					continue
				}

				log.Println("READ - REF: ", string(buf))
			}
		}
	}()
}
