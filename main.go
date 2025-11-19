package main

import (
	"log"
	"github.com/joho/godotenv"

	overview "github.com/pachoR/stonks-data-files/overview"
)

func init () {
	godotenv.Load()
}

func main() {

	err := overview.CreateOverviewIndex()
	if err != nil {
		log.Fatalf("Err! Couldn't create overview index: %s", err.Error())
	}
	// overview.CreateOverviewDataFile()
	err = overview.IngestOverviewData()
	if err != nil {
		log.Fatalf("Err! Ingesting process failed: %s", err.Error())
	}
}
