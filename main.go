package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

var ENDPOINT = "https://www.rescuetime.com/anapi/data"

func main() {
	url, err := url.Parse(ENDPOINT)
	if err != nil {
		log.Fatal("Could not parse URL")
	}
	query := url.Query()
	query.Add("key", "B63EI72gQZsau2QXpFX3UGqGSfVNSJfotneLLtJP")
	query.Add("format", "csv")
	query.Add("pv", "interval")
	query.Add("rs", "day")
	query.Add("rb", "2023-02-01")
	query.Add("re", "2023-02-28")
	query.Add("rk", "productivity")

	url.RawQuery = query.Encode()

	res, err := http.Get(url.String())
	if err != nil {
		log.Fatalf("Could not get url: %s", url)
	}

	defer res.Body.Close()

	csvReader := csv.NewReader(res.Body)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Could not read CSV data")
	}

	rows := readData(data)

	fmt.Printf("Data: %+v\n", rows)
}
