package orderbook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type OrderBook struct {
	// Name of the exchanger
	Exchanger string
	Bids      []*Order
	Asks      []*Order
}

type Order struct {
	Price     float64
	Volume    float64
	Timestamp float64
}

// NewOrderbook returns a new orderbook. An error is returned when orders are
// not sorted as expected or when they are no orders in at least one side of the book.
func NewOrderbook(Exchanger string, bids, asks []*Order) (*OrderBook, error) {
	// verify orderbook isn't empty
	if len(bids) == 0 {
		return nil, fmt.Errorf("Orderbook: no ask orders (%s).", Exchanger)
	} else if len(asks) == 0 {
		return nil, fmt.Errorf("Orderbook: no bid orders (%s).", Exchanger)
	}

	// verify bid orders are sorted
	maxBid := bids[0].Price
	for _, o := range bids[1:] {
		if o.Price > maxBid {
			return nil, fmt.Errorf("Orderbook: %s bid orders are not sorted.", Exchanger)
		}
		maxBid = o.Price
	}

	// verify ask orders are sorted
	minAsk := asks[0].Price
	for _, o := range asks[1:] {
		if o.Price < minAsk {
			return nil, fmt.Errorf("Orderbook: %s ask orders are not sorted", Exchanger)
		}
		minAsk = o.Price
	}

	return &OrderBook{Exchanger, bids, asks}, nil
}

func FetchOrderBook(url string, v interface{}) error {
	// create the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}

	// execute the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}

	// read the response body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	return json.Unmarshal(body, v)
}
