package main

import (
	"go_practicum/app/config"
	"go_practicum/app/shortener"
	"log"
)

func main() {
	cfg := config.NewConfig()

	if err := shortener.Start(*cfg); err != nil {
		log.Fatal(err)
	}
}
