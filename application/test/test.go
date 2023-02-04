package main

import (
	"log"
	"time"
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

	log.Println(time.Now().UTC())
}
