package main

import "strconv"

type Productivity = map[int]int

const (
	VeryProductive  int = 2
	Productive          = 1
	Neutral             = 0
	Distracting         = -1
	VeryDistracting     = -2
)

func readData(data [][]string) map[string]Productivity {
	list := make(map[string]Productivity, len(data)-1)
	for i, line := range data {
		if i == 0 {
			continue
		}

		date := line[0]
		if list[date] == nil {
			list[date] = make(Productivity, 5)
		}

		index, err := strconv.Atoi(line[3])
		if err != nil {
			panic(err)
		}

		list[date][index], err = strconv.Atoi(line[1])
		if err != nil {
			panic(err)
		}
	}
	return list
}
