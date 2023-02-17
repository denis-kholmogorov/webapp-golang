package main

import (
	"log"
	"math"
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
	//tomorrow := NextDay(today

	//posts := domain.Posts{TotalElement: 8}
	//i := posts.TotalPages / 3

	log.Println(math.Ceil(1.0 / 3))
}
