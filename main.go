package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

var ENDPOINT = "https://www.rescuetime.com/anapi/data"

func center(w int, s string) {
	fmt.Printf("%*s", -w, fmt.Sprintf("%*s", (w+len(s))/2, s))
}

func printRow(key string, data Productivity) {
	center(10, key)
	for _, v := range data {
		center(10, fmt.Sprint(v))
	}
	fmt.Println()
}

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

	for k, i := range rows {
		printRow(k, i)
	}
}
