package main

import (
	"fmt"
	"go_practicum/app/config"
	"go_practicum/app/shortener"
	"log"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	typeGlVars()

	cfg := config.NewConfig()

	if err := shortener.Start(*cfg); err != nil {
		log.Fatal(err)
	}
}

func typeGlVars() {
	fmt.Printf("Build version: %v\n", buildVersion)
	fmt.Printf("Build date: %v\n", buildDate)
	fmt.Printf("Build commit: %v\n", buildCommit)
}
