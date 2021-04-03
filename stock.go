package main

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type stock struct {
	Symbol   string
	Quantity int
	Price    quote
}

type quote struct {
	CurrPrice float64
	PrevClose float64
}

func (s *stock) getPrice() error {
	if err := s.Price.update(s.Symbol); err != nil {
		return err
	}

	return nil
}

func (q *quote) fetch(url string) ([]byte, error) {
	// http client with timeout
	client := http.Client{Timeout: 10 * time.Second}

	// fetch URL
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// convert to byte slice
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (q *quote) update(symbol string) error {
	resp, err := q.fetch("https://query2.finance.yahoo.com/v10/finance/quoteSummary/" +
		symbol + "?formatted=false&modules=price")
	if err != nil {
		return err
	}

	// get current stock price
	re := regexp.MustCompile(`regularMarketPrice\":[0-9]*\.[0-9]+`)
	currPrice := string(re.Find(resp))
	currPrice = strings.TrimPrefix(currPrice, "regularMarketPrice\":")
	q.CurrPrice, err = strconv.ParseFloat(currPrice, 64)
	if err != nil {
		return err
	}

	// get previous close
	re = regexp.MustCompile(`regularMarketPreviousClose\":[0-9]*\.[0-9]+`)
	prevClose := string(re.Find(resp))
	prevClose = strings.TrimPrefix(prevClose, "regularMarketPreviousClose\":")
	q.PrevClose, err = strconv.ParseFloat(prevClose, 64)
	if err != nil {
		return err
	}

	return nil
}
