package main

import "strconv"

type Productivity struct {
	Time         int
	Productivity int
}

const (
	VeryProductive  int = 2
	Productive          = 1
	Neutral             = 0
	Distracting         = -1
	VeryDistracting     = -2
)

func readData(data [][]string) map[string][]Productivity {
	list := make(map[string][]Productivity, len(data)-1)
	var err error
	for i, line := range data {
		if i == 0 {
			continue
		}
		var date string
		var row Productivity
		for j, field := range line {
			switch j {
			case 0:
				date = field
			case 1:
				row.Time, err = strconv.Atoi(field)
				if err != nil {
					panic(err)
				}
			case 2:
				continue
			case 3:
				row.Productivity, err = strconv.Atoi(field)
				if err != nil {
					panic(err)
				}
			}
		}
		if list[date] == nil {
			list[date] = make([]Productivity, 5)
		}
		list[date] = append(list[date], row)
	}
	return list
}
