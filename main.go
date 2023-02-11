package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

var ENDPOINT = "https://www.rescuetime.com/anapi/data"
var COLUMNS = [5]int{VeryProductive, Productive, Neutral, Distracting, VeryDistracting}

func center(w int, s string) {
	fmt.Printf("%*s", -w, fmt.Sprintf("%*s", (w+len(s))/2, s))
}

func printRow(key time.Time, data Productivity) {
	center(10, key.Format("2006-01-02"))
	for _, i := range COLUMNS {
		item := data[i]
		if item == 0 {
			center(10, "00:00")
			continue
		}
		hours := data[i] / 3600
		minutes := (data[i] % 3600) / 60
		center(10, fmt.Sprintf("%02d:%02d", hours, minutes))
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
