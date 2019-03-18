package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.Now())
	fmt.Println(time.UTC)
	t := time.Date(2019, 2, 7, 0, 0, 0, 0, time.UTC)
	fmt.Println(t.Format("2006-01-02"))
	fmt.Println(t.AddDate(0, 0, -1).Format("2006-01-02"))
	fmt.Println(t.Weekday())
	fmt.Println(int(t.Weekday()))
}
