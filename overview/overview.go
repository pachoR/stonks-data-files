package overview

import (
	"bufio"
	"context"
	"encoding/json"
	"strings"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	http "github.com/pachoR/go-libs/http"
	oslib "github.com/pachoR/go-libs/oslib"
	postgres "github.com/pachoR/go-libs/postgreslib"
	"github.com/pachoR/stonks-data-files/overview/types"
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
	url := fmt.Sprintf("%s/search?q=%s&exchange=US", os.Getenv("FINN_URL"), symbol)

	header := map[string]string {
		"X-Finnhub-Token": os.Getenv("FINN_KEY"),
	}

	res, err := http.GetWithHeader(url, header)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	bytes, _ := io.ReadAll(res.Body)
	return bytes, err
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
	var notFoundSymbols []string

	for _, symbol := range symbols {
		startFetch := time.Now()

		fetchedSymbolOverview, err := getOverviewData(symbol)
		if err != nil {
			return fmt.Errorf("Error fetchingSymbol %s: %s", symbol, err.Error())
		}

		var overviewSymbolMeta types.SymbolOverviewMeta
		err = json.Unmarshal(fetchedSymbolOverview, &overviewSymbolMeta)

		var overviewSymbolDetail types.SymbolOverview
		for _, symDetail := range overviewSymbolMeta.Result {
			if strings.EqualFold(symDetail.Symbol, symbol) {
				overviewSymbolDetail = symDetail
			}
		}

		if len(overviewSymbolDetail.Symbol) == 0 {
			notFoundSymbols = append(notFoundSymbols, symbol)
			continue
		}

		trimmedJson, _ := json.Marshal(overviewSymbolDetail)
		_, err = f.WriteString(string(trimmedJson) + "\n")
		if err != nil {
			log.Printf("Error on writing the following json object: %s", string(trimmedJson))
			return fmt.Errorf("Error writing symbol %s: %s", symbol, err.Error())
		}

		gap := time.Second - time.Since(startFetch)
		if gap > 0 {
			time.Sleep(gap)
		} else {
			gap = 0
		}
		log.Printf("Writing info for symbol: %s\nWAIT for %v", symbol, gap)
	}

	fmt.Printf("Not found symbols: %d\n", len(notFoundSymbols))
	for _, sym := range notFoundSymbols {
		fmt.Println(sym)
	}

	return nil
}

func CreateOverviewIndex() error {
	err := oslib.CreateIndex(IndexName, "./mapping/overview.json")
	if err != nil {
		return err
	}

	return nil
}

func IngestOverviewData() error {
	f, err := os.Open("local/overview.ndjson")
	if err != nil {
		return fmt.Errorf("Err! Couldn't open local/overview.ndjson file %s", err.Error())
	}

	fScanner := bufio.NewScanner(f)
	fScanner.Split(bufio.ScanLines)

	for fScanner.Scan() {
		ovwBytes := []byte(fScanner.Text())
		var overview types.SymbolOverview
		err := json.Unmarshal(ovwBytes, &overview)
		if err != nil {
			fmt.Printf("ERR! Counlnt unmarshal: %s\n%s\n", err, string(ovwBytes))
			continue
		}

		err = oslib.IngestFromStruct(IndexName, overview)
		if err != nil {
			fmt.Printf("ERR! error ingesting from struct: %s\n", err.Error())
		}
	}
	return nil
}