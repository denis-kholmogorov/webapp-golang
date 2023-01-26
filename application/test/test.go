package main

import (
	"fmt"
	"math/big"
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
	b := big.NewInt(1)
	fmt.Println(b)
}
