package main

import (
	"github.com/joho/godotenv"
	overview "github.com/pachoR/stonks-data-files/overview"
)

func init () {
	godotenv.Load()
}

func main() {
	overview.CreateOverviewDataFile()
}
