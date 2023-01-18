package main

import (
	"fmt"
	"strings"
)

func main() {
	rowtoken := "Bearer asdasdasdasd"
	before := strings.TrimPrefix(rowtoken, "Bearer ")
	fmt.Println(before)
}
