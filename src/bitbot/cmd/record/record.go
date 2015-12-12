package main

import (
	"flag"
	"log"
	"time"

	"bitbot/database"
	"bitbot/exchanger/bitfinex"
	"bitbot/exchanger/btce"
	"bitbot/exchanger/cex"
	"bitbot/exchanger/hitbtc"
	"bitbot/exchanger/kraken"
	"bitbot/exchanger/orderbook"
)

// TODO: this code is not panic safe
var (
	dbPath      = flag.String("d", "./data/book.sql", "SQLite database path.")
	periodicity = flag.Int64("t", 5, "Periodicity expressed in seconds.")
)

type exchanger struct {
	name string
	f    func(string) (*orderbook.OrderBook, error)
}

func main() {
	log.Println("Start recording...")

	flag.Parse()
	pair := "BTC_USD"

	db := database.Open(*dbPath)
	defer db.Close()

	exchangers := []*exchanger{
		&exchanger{"hitbtc", hitbtc.OrderBook},
		&exchanger{"bitfinex", bitfinex.OrderBook},
		&exchanger{"btce", btce.OrderBook},
		&exchanger{"kraken", kraken.OrderBook},
		&exchanger{"cex", cex.OrderBook},
	}

	database.CreateTable(db, pair)

	for {
		for _, e := range exchangers {
			// TODO: timeout after 2 sec
			go func(e *exchanger) {
				if r := fetchRecord(e, pair); r != nil {
					database.SaveRecord(db, pair, r)
				}
			}(e)
		}

		time.Sleep(time.Duration(*periodicity) * time.Second)
	}
}

func fetchRecord(e *exchanger, pair string) *database.Record {
	log.Printf("Fetching %s...", e.name)
	start := time.Now().UnixNano()
	book, err := e.f(pair)
	end := time.Now().UnixNano()

	if err != nil {
		log.Println(err)
		return nil
	}

	return &database.Record{
		Exchanger: e.name,
		StartTime: start,
		EndTime:   end,
		Bids:      book.Bids,
		Asks:      book.Asks,
	}
}
