package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/7phs/area-51/app/config"
	data_stream "github.com/7phs/area-51/app/data-stream"
	"github.com/7phs/area-51/app/watcher"
)

func main() {
	conf := config.Parse()
	if err := conf.Validate(); err != nil {
		log.Fatal("invalid config: ", err)
	}

	w, err := watcher.NewWatcher()
	if err != nil {
		log.Fatal("failed to init watcher: ", err)
	}

	referenceQueue, err := w.WatchFileChanges(conf.ReferenceFile)
	if err != nil {
		log.Fatal("failed to init reference queue changes: ", err)
	}

	referenceStream, err := data_stream.NewDataStream(referenceQueue)
	if err != nil {
		log.Fatal("failed to init reference data stream: ", err)
	}

	rawQueue, err := w.WatchFileChanges(conf.RawDataFile)
	if err != nil {
		log.Fatal("failed to init raw queue changes: ", err)
	}

	rawStream, err := data_stream.NewDataStream(rawQueue)
	if err != nil {
		log.Fatal("failed to init raw data stream: ", err)
	}

	referenceStream.Start()
	rawStream.Start()
	w.Start()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		log.Println("interrupt")

		cancel()
	}()

	<-ctx.Done()

	w.Stop()
	referenceStream.Stop()
	rawStream.Stop()

	log.Println("success")
}
