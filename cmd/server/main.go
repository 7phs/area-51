package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/7phs/area-51/app/config"
	"github.com/7phs/area-51/app/server"
)

func main() {
	conf := config.Parse()
	if err := conf.Validate(); err != nil {
		log.Fatal("invalid config: ", err)
	}

	log.Println("init")

	srv, err := server.New(conf)
	if err != nil {
		log.Fatal("failed to init server: ", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		log.Println("interrupt")

		cancel()
	}()

	log.Println("start")

	srv.Start()

	<-ctx.Done()

	log.Println("stopping...")

	srv.Stop()

	log.Println("stop")
}
