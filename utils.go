package main

import (
	"sort"
	"time"
)

func SortDescending(examples *[]Example) {
	sort.Slice(*examples, func(i, j int) bool {
		date1, _ := time.Parse("2-Jan-2006", (*examples)[i].Date)
		date2, _ := time.Parse("2-Jan-2006", (*examples)[j].Date)
		if date1.After(date2) {
			return true
		} else {
			return false
		}
	})
}
