package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"sort"
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

func makeURL() (*url.URL, float64, float64) {
	now := time.Now()
	startOfMonth := now.Add(-time.Hour * time.Duration(24*(now.Day()-1)))

	key := flag.String("key", "", "RescueTime API key")
	start := flag.String("start", startOfMonth.Format("2006-01-02"), "Sets start date for fetch period")
	end := flag.String("end", now.Format("2006-01-02"), "Sets end date for fetch period")
	wage := flag.Float64("wage", 375, "Amount of money you earn per productive hour")
	multiplier := flag.Float64("multiplier", 1, "Multiply all productive time")

	flag.Parse()

	if *key == "" {
		var file *os.File
		var err error
		switch runtime.GOOS {
		case "windows":
			file, err = os.Open(os.ExpandEnv("%LOCALAPPDATA%/RescueTime.com/rescuetimed.json"))
		case "linux":
			file, err = os.Open(os.ExpandEnv("$HOME/.config/RescueTime.com/rescuetimed.json"))
		default:
			file, err = os.Open(os.ExpandEnv("$HOME/Library/RescueTime.com/rescuetimed.json"))
		}
		if err != nil {
			log.Fatal("No RescueTime API key is provided")
		}
		bytes, err := ioutil.ReadAll(file)
		var config RescueTimeConfig
		json.Unmarshal(bytes, &config)
		*key = config.Key
		if *key == "" {
			log.Fatal("No RescueTime API key is provided")
		}
	}

	url, err := url.Parse(ENDPOINT)
	if err != nil {
		log.Fatal("Could not parse URL")
	}
	query := url.Query()
	query.Add("key", *key)
	query.Add("format", "csv")
	query.Add("pv", "interval")
	query.Add("rs", "day")
	query.Add("rb", *start)
	query.Add("re", *end)
	query.Add("rk", "productivity")

	url.RawQuery = query.Encode()
	return url, *wage, *multiplier
}

func printSummary(data map[time.Time]Productivity, wage float64, multiplier float64) {
	sum := float64(0)
	for _, v := range data {
		for c := range COLUMNS {
			// do not use seconds in calculations
			sum += float64(v[c] - (v[c] % 60))
		}
	}

	productiveHours := (sum / 3600 * multiplier)
	fmt.Print(color.Ize(color.Green, "Total productive hours: "))
	fmt.Printf("%.2f\n", productiveHours)

	fmt.Print(color.Ize(color.Green, "Your wage is "))
	fmt.Printf("%.2f\n", productiveHours*wage)
}

func main() {
	url, wage, multiplier := makeURL()

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

	keys := make([]time.Time, 0, len(rows))
	for t := range rows {
		keys = append(keys, t)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Before(keys[j])
	})

	printHead()
	for _, k := range keys {
		printRow(k, rows[k])
	}
	printSummary(rows, wage, multiplier)
}
