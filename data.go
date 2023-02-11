package main

import "strconv"

type RescueData struct {
	Date         string
	Time         int
	People       int
	Productivity int
}

const (
	VeryProductive  int = 2
	Productive          = 1
	Neutral             = 0
	Distracting         = -1
	VeryDistracting     = -2
)

func readData(data [][]string) []RescueData {
	var list []RescueData
	var err error
	for i, line := range data {
		if i == 0 {
			continue
		}
		var row RescueData
		for j, field := range line {
			switch j {
			case 0:
				row.Date = field
			case 1:
				row.Time, err = strconv.Atoi(field)
				if err != nil {
					panic(err)
				}
			case 2:
				row.People, err = strconv.Atoi(field)
				if err != nil {
					panic(err)
				}
			case 3:
				row.Productivity, err = strconv.Atoi(field)
				if err != nil {
					panic(err)
				}
			}
		}
		list = append(list, row)
	}
	return list
}
