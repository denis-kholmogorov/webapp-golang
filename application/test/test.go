package main

import (
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

	//i := int64(1646915447)

	//m := map[string]string{"Четвертый": " пост"}
	//for s, s2 := range m {
	//	fmt.Println(s + s2)
	//}
	//now := time.Now()
	//t := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	//var today Weekday = Sunday
	//tomorrow := NextDay(today

	//posts := domain.Posts{TotalElement: 8}
	//i := posts.TotalPages / 3

	//fmt.Println(time.Unix(i, 0).Format("2006-01-02T03:04:05.999999999Z"))
	parse, err := time.Parse("2006-01-02T15:04:05.999999999Z", "1970-08-01T16:14:54.240Z")
	if err != nil {
		return
	}
	println(parse.String())
}
