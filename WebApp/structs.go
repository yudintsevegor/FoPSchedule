package main

import (
	"database/sql"
	"fopSchedule/master/common"
	"net/http"
	"sync"
	"time"
)

type SubjectsInfo struct {
	Subject       common.Subject
	Number        int
	IsOdd         bool
	IsAllDay      bool
	LessonStartAt string
	TimeNow       time.Time
	SemesterEnd   string
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
	Mutex    *sync.Mutex
	DB       *sql.DB
}

type ServerError struct {
	Error string
}
