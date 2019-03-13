package main

import (
	"regexp"
	"time"
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

type DataToParsingLine struct {
	Departments      []Department
	AllGroups        []string
	ResultFromReqexp []string
	InsertedGroups   []string
	Lesson           Subject
	RegexpInterval   *regexp.Regexp
}

type Interval struct {
	Start int
	End   int
}

type LessonRange struct {
	Start string
	End   string
}

type DataToParsingAt struct {
	Lesson      Subject
	Number      int
	Parity      bool
	IsAllDay    bool
	StartTime   string
	Time        time.Time
	SemesterEnd string
}

var (
	columns = " ( first, second, third, fourth, fifth ) "
	quesStr = " ( ?, ?, ?, ?, ? ) "

	reUpp  = regexp.MustCompile("([А-Я]){5,}")
	rePerc = regexp.MustCompile("(.*)%(.*)%(.*)")
	reAt   = regexp.MustCompile("(.*)@(.*)")
	reNum  = regexp.MustCompile(`([0-9]+М*)`)

	practice = "Преддипломная практика"
	WAR      = "ВОЕННАЯ ПОДГОТОВКА"
	war      = "Военная подготовка"
	MFK      = "МЕЖФАКУЛЬТЕТСКИЕ КУРСЫ"
	MFKabbr = "МФК"
	mfk = "Межфакультетские курсы"
	prac = "Общий физический практикум"
	specprac = "Специальный физический практикум"
	prac201 = "Практикум по радиоэлектронике"
	phys = "Физическая культура"
	research = "Научно-исследовательская практика"
	astroProblems = "Современные проблемы астрономии"
	NIS = "НИС"

	cases = WAR + " " + war + " " + MFK + " " + mfk + " " + MFKabbr
	astr = "астр."
	moscowTime    = "+03:00"
	timeIntervals = map[int]LessonRange{
		0: {Start: "T9:00:00" + moscowTime, End: "T10:35:00" + moscowTime},
		1: {Start: "T10:50:00" + moscowTime, End: "T12:25:00" + moscowTime},
		2: {Start: "T13:30:00" + moscowTime, End: "T15:05:00" + moscowTime},
		3: {Start: "T15:20:00" + moscowTime, End: "T16:55:00" + moscowTime},
		4: {Start: "T17:05:00" + moscowTime, End: "T18:40:00" + moscowTime},
	}
)
