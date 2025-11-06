package overview

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	http "github.com/pachoR/go-libs/http"
	postgres "github.com/pachoR/go-libs/postgreslib"
)

var symbols []string

func getAllSymbols() error {
	conn, err := postgres.GetConnection()
	if err != nil {
		log.Printf("Error getting connection on the database: %s\n", err.Error())
		return err
	}

	rows, err := conn.Query(context.Background(), "SELECT symbol FROM symbols")
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var sym string
		err = rows.Scan(&sym)
		if  err != nil {
			return fmt.Errorf("Error scanning symbols at getAllSymbols: %s", err.Error())
		}
		symbols = append(symbols, sym)
	}

	err = rows.Err()
	if err != nil {
		return fmt.Errorf("Error reading rows: %s", err.Error())
	}

	totalSymbols := len(symbols)
	log.Printf("%d symbols retrived", totalSymbols)
	if totalSymbols > 0 {
		log.Printf("[0]: %s", symbols[0])
		log.Printf("[%d]: %s", totalSymbols-1, symbols[totalSymbols-1])
	}
	return nil
}

func getOverviewData(symbol string) ([]byte, error) {
	url := fmt.Sprintf("%s&symbol=%s&apikey=%s", getOverviewURL(), symbol, getApiKey())

	bytes, err := http.GetBody(url)
	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}

func timer() func() {
	start := time.Now()
	return func() {
		log.Printf("CreateOverviewDataFile took %v\n", time.Since(start))
	}
}

func CreateOverviewDataFile() error {
	defer timer()()

	err := getAllSymbols()
	if err != nil {
		return fmt.Errorf("Error at getAllSymbols: %s", err.Error())
	}

	if len(symbols) == 0{
		return fmt.Errorf("SYMBOLS QUERY RETURNED EMPTY")
	}

	f, err := os.Create("local/overview.ndjson")
	if err != nil {
		return fmt.Errorf("Error creating file at local/overview.json: %s", err.Error())
	}
	defer f.Close()

	for _, symbol := range symbols {
		fetchedSymbolOverview, err := getOverviewData(symbol)
		if err != nil {
			return fmt.Errorf("Error fetchingSymbol %s: %s", symbol, err.Error())
		}

		var jsonSymbol interface{}
		err = json.Unmarshal(fetchedSymbolOverview, &jsonSymbol)

		trimmedJson, _ := json.Marshal(jsonSymbol)
		_, err = f.WriteString(string(trimmedJson) + "\n")
		if err != nil {
			log.Printf("Error on writing the following json object: %s", string(trimmedJson))
			return fmt.Errorf("Error writing symbol %s: %s", symbol, err.Error())
		}
		log.Printf("Writing info for symbol: %s", symbol)
	}

	return nil
}
