package main

import (
	"fopSchedule/common"
	"regexp"
	"time"
)

type Department struct {
	Number  string
	Lessons []common.Subject
}

type DataToParsingLine struct {
	Departments      []Department
	AllGroups        []string
	ResultFromReqexp []string
	InsertedGroups   []string
	Lesson           common.Subject
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
	Lesson      common.Subject
	Number      int
	Parity      bool
	IsAllDay    bool
	StartTime   string
	Time        time.Time
	SemesterEnd string
}

type Template struct {
	Course string
	Group  string
}
