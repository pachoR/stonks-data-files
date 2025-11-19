package overview

import (
	"fmt"
	"os"
	"sync"
)

const IndexName string = "symbols-overview"

var (
	overviewURL string
	overviewURLOnce sync.Once
)

var (
	apiKey string
	apiKeyOnce sync.Once
)

func getOverviewURL() string {
	overviewURLOnce.Do(func () {
		overviewURL = fmt.Sprintf("%sfunction=OVERVIEW", os.Getenv("ALPHA_URL"))
	})
	return overviewURL
}

func getApiKey() string {
	apiKeyOnce.Do(func () {
		apiKey = fmt.Sprintf("%s", os.Getenv("ALPHA_KEY"))
	})
	return apiKey
}
