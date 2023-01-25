package main

import (
	"fmt"
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
	fmt.Println("today =", Sunday, "tomorrow =", Monday)
	capd := time.Now()
	after := time.Time{}.Before(capd.Add(time.Minute * 2))
	fmt.Println(after)
}
