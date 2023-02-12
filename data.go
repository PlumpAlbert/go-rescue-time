package main

import (
	"strconv"
	"time"
)

type Productivity = map[int]int

type RescueTimeConfig struct {
	Key string `json:"data_key"`
}

const (
	VeryProductive  int = 2
	Productive          = 1
	Neutral             = 0
	Distracting         = -1
	VeryDistracting     = -2
)

func readData(data [][]string) map[time.Time]Productivity {
	list := make(map[time.Time]Productivity, 1)
	for i, line := range data {
		if i == 0 {
			continue
		}

		date, err := time.Parse("2006-01-02T15:04:05", line[0])
		if err != nil {
			panic(err)
		}
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
