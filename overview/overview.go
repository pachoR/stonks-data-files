package overview

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	postgres "github.com/pachoR/go-libs/postgreslib"
)

var (
	overviewURL string
	overviewURLOnce sync.Once
)

var symbols []string

func getOverviewURL() string {
	overviewURLOnce.Do(func () {
		overviewURL = fmt.Sprintf("%sfunction=OVERVIEW", os.Getenv("ALPHA_URL"))
	})
	return overviewURL
}

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

//func GetOverviewData(symbol string) (SymbolOverview error) {
//	bytes, err := http.GetBody(getOverviewURL())
//	if err != nil {
//		return nil, err
//	}
//}

func CreateOverviewDataFile() error {
	err := getAllSymbols()
	if err != nil {
		log.Printf("Error at getAllSymbols: %s", err.Error())
		return err
	}

	return nil
}
