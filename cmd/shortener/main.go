package main

import (
	"log"

	"go_practicum/internal/app/shortener"
)

func main() {
	config := shortener.NewConfig() //getting config with default values
	config.ConfigureViaEnv()        //overwritting config using values from env
	config.ConfigureViaFlags()      //overwritting config using values from flags

	if err := shortener.Start(config); err != nil {
		log.Fatal(err)
	}
}
