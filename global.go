package main

import (
	"regexp"
)

type Subject struct {
	Name   string
	Lector string
	Room   string
}

type Department struct {
	Number  string
	Lessons []Subject
}

type Interval struct {
	Start int
	End   int
}

type LessonRange struct {
	Start string
	End   string
}

var (
	reUpp  = regexp.MustCompile("([А-Я]){2,}")
	rePerc = regexp.MustCompile("(.*)%(.*)%(.*)")
	reNum  = regexp.MustCompile(`([0-9]+)`)
	reAt   = regexp.MustCompile("(.*)@(.*)")

	practice = "Преддипломная практика"
	war      = "ВОЕННАЯ ПОДГОТОВКА"
	mfk      = "МЕЖФАКУЛЬТЕТСКИЕ КУРСЫ"

	moscowTime    = "+03:00"
	timeIntervals = map[int]LessonRange{
		0: {Start: "T9:00:00" + moscowTime, End: "T10:35:00" + moscowTime},
		1: {Start: "T10:50:00" + moscowTime, End: "T12:25:00" + moscowTime},
		2: {Start: "T13:30:00" + moscowTime, End: "T15:05:00" + moscowTime},
		3: {Start: "T15:20:00" + moscowTime, End: "T16:55:00" + moscowTime},
		4: {Start: "T17:05:00" + moscowTime, End: "T18:40:00" + moscowTime},
	}
)
