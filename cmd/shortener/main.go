package main

import (
	"context"
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
	ctx := context.Background()

	if err := shortener.Start(ctx, *cfg); err != nil {
		log.Fatal(err)
	}
}

func typeGlVars() {
	fmt.Printf("Build version: %v\n", buildVersion)
	fmt.Printf("Build date: %v\n", buildDate)
	fmt.Printf("Build commit: %v\n", buildCommit)
}
