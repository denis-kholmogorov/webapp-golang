package main

import (
	"log"
	"strconv"
)

type Weekday int

const (
	Monday Weekday = iota + 1
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

func NextDay(day Weekday) Weekday {
	return (day % 7) + 1
}

func main() {
	//var today Weekday = Sunday
	//tomorrow := NextDay(today)
	var f float64 = 1.2654
	float := strconv.FormatFloat(f, 'G', -1, 64)
	log.Println("F = " + float)
}
