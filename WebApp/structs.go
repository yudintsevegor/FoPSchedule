package main

import (
	"net/http"
	"regexp"
	"time"
)

type Subject struct {
	Name   string
	Lector string
	Room   string
	Parity string
}

type Department struct {
	Number  string
	Lessons []Subject
}

type DataToParsingLine struct {
	Departments []Department
	AllGroups   []string

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

type Template struct {
	Course string
	Group  string
}

type User struct {
	Client     *http.Client
	Email      string
	PathAction string
}

type UserInfo struct {
	Email string `json:"email"`
}

type Handler struct {
	Sessions map[string]User
}

type ServerError struct {
	Error string
}
