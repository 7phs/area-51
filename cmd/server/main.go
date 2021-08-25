package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/7phs/area-51/app/config"
	"github.com/7phs/area-51/app/server"
)

func main() {
	conf := config.Parse()
	if err := conf.Validate(); err != nil {
		log.Fatal(time.Now(), "invalid config: ", err)
	}

	log.Println(time.Now(), "init")

	srv, err := server.New(conf)
	if err != nil {
		log.Fatal(time.Now(), "failed to init server: ", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		log.Println(time.Now(), "interrupt")

		cancel()
	}()

	log.Println(time.Now(), "start")

	srv.Start()

	<-ctx.Done()

	log.Println(time.Now(), "stopping...")

	srv.Stop()

	log.Println(time.Now(), "stop")
}
