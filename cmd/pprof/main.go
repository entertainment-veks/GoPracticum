package main

import (
	"go_practicum/app/config"
	"go_practicum/app/pprof"
	"go_practicum/app/shortener"
	"log"
	_ "net/http/pprof"
	"time"
)

func main() {
	cfg := config.NewConfig()

	go func() {
		time.Sleep(time.Second)
		go pprof.Loading(*cfg)
	}()

	if err := shortener.Start(*cfg); err != nil {
		log.Fatal(err)
	}
}
