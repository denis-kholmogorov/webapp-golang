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
	now := time.Now()
	t := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	//var today Weekday = Sunday
	//tomorrow := NextDay(today

	//posts := domain.Posts{TotalElement: 8}
	//i := posts.TotalPages / 3

	fmt.Println(t.Format("2006-01-02T03:04:05Z"))
}
