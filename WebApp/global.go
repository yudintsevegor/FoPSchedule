package main

import (
	"regexp"
)

const (
	maxConnections = 100

	cookieURL  = "/cookie"
	cookieName = "fopschedule"

	mainHTMLPage  = "mainPage.html"
	indexHTMLPage = "index.html"

	calendarName = "Shedule"
	urlCalendar  = "https://calendar.google.com"
	timeLayout   = "2006-01-02"

	columns = " ( first, second, third, fourth, fifth ) "
	quesStr = " ( ?, ?, ?, ?, ? ) "

	practice      = "Преддипломная практика"
	WAR           = "ВОЕННАЯ ПОДГОТОВКА"
	war           = "Военная подготовка"
	MFK           = "МЕЖФАКУЛЬТЕТСКИЕ КУРСЫ"
	MFKabbr       = "МФК"
	mfk           = "Межфакультетские курсы"
	prac          = "Общий физический практикум"
	specprac      = "Специальный физический практикум"
	prac201       = "Практикум по радиоэлектронике"
	phys          = "Физическая культура"
	research      = "Научно-исследовательская практика"
	astroProblems = "Современные проблемы астрономии"
	NIS           = "НИС"
	astr          = "астр."

	lessonCases = WAR + " " + war + " " + MFK + " " + mfk + " " + MFKabbr
)

var (
	// to set another color for event if event is for all groups
	reUpp = regexp.MustCompile("([А-Я]){5,}")

	subGroups = map[string][]string{
		"341":  []string{"341а", "341б"},
		"441":  []string{"441а", "441б"},
		"141М": []string{"141Ма", "141Мб"},
		"241М": []string{"241Ма", "241Мб"},
		"316":  []string{"316а", "316б"},
		"416":  []string{"416а", "416б"},
		"116М": []string{"116Ма", "116Мб"},
		"216М": []string{"216Ма", "216Мб"},
	}

	south = map[string]string{
		"5-23": "",
		"5-24": "",
		"5-25": "",
		"5-26": "",
		"5-27": "",
		"5-38": "",
		"5-39": "",
		"5-40": "",
		"5-41": "",
		"5-42": "",
		"5-18": "",
		"5-19": "",
	}
	north = map[string]string{
		"5-33":   "",
		"5-34":   "",
		"5-35":   "",
		"5-36":   "",
		"5-37":   "",
		"5-44":   "",
		"5-45":   "",
		"5-46":   "",
		"5-47":   "",
		"5-48":   "",
		"5-49":   "",
		"5-50":   "",
		"5-51":   "",
		"5-52":   "",
		"5-53":   "",
		"5-61":   "",
		"5-62":   "",
		"5-68":   "",
		"Л.каб.": "",
	}

	audience = map[string]string{
		"СФА":              "",
		"ЮФА":              "",
		"ЦФА":              "",
		"Ауд. им. Хохлова": "",
	}

	moscowTime    = "+03:00"
	timeIntervals = map[int]LessonRange{
		0: {Start: "T9:00:00" + moscowTime, End: "T10:35:00" + moscowTime},
		1: {Start: "T10:50:00" + moscowTime, End: "T12:25:00" + moscowTime},
		2: {Start: "T13:30:00" + moscowTime, End: "T15:05:00" + moscowTime},
		3: {Start: "T15:20:00" + moscowTime, End: "T16:55:00" + moscowTime},
		4: {Start: "T17:05:00" + moscowTime, End: "T18:40:00" + moscowTime},
	}
)
