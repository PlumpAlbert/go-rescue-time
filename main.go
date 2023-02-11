package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/TwiN/go-color"
)

var ENDPOINT = "https://www.rescuetime.com/anapi/data"
var COLUMNS = [5]int{VeryProductive, Productive, Neutral, Distracting, VeryDistracting}

func center(w int, s string) string {
	return fmt.Sprintf("%*s", -w, fmt.Sprintf("%*s", (w+len(s))/2, s))
}

func printHead() {
	fmt.Printf("%*s", -9, fmt.Sprintf("%*s", 5, "📅"))
	fmt.Printf("%*s", -9, fmt.Sprintf("%*s", 5, "🥰"))
	fmt.Printf("%*s", -9, fmt.Sprintf("%*s", 5, "😊"))
	fmt.Printf("%*s", -9, fmt.Sprintf("%*s", 5, "😐"))
	fmt.Printf("%*s", -9, fmt.Sprintf("%*s", 5, "😟"))
	fmt.Printf("%*s", -9, fmt.Sprintf("%*s", 5, "😢"))
	fmt.Println()
}

func printRow(key time.Time, data Productivity) {
	fmt.Print(color.Ize(color.Green, center(10, key.Format("2006-01-02"))))
	for _, i := range COLUMNS {
		var s string
		item := data[i]
		if item == 0 {
			s = center(10, "00:00")
		} else {
			hours := data[i] / 3600
			minutes := (data[i] % 3600) / 60
			s = center(10, fmt.Sprintf("%02d:%02d", hours, minutes))
		}

		var c string
		switch i {
		case VeryProductive:
			c = color.Blue
		case Productive:
			c = color.Cyan
		case Neutral:
			c = color.Gray
		case Distracting:
			c = color.Yellow
		case VeryDistracting:
			c = color.Red
		}
		fmt.Print(color.Ize(c, s))
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

	printHead()
	for k, i := range rows {
		printRow(k, i)
	}
}
