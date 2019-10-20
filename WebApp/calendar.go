package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"strings"
	"time"

	"google.golang.org/api/calendar/v3"

	_ "github.com/go-sql-driver/mysql"
)

const (
	endSummerSemester = "0601"
	endWinterSemester = "1231"
)

func (u *User) putData(db *sql.DB, client *http.Client, group string) error {
	clndr := &calendar.Calendar{
		Summary: calendarName + group,
	}

	srvc, err := calendar.New(client)
	if err != nil {
		return err
	}

	insertedCalendar, err := srvc.Calendars.Insert(clndr).Do()
	if err != nil {
		return err
	}

	var (
		isOdd       bool
		day         int
		endSemester string
	)

	timeNow := time.Now()
	year := timeNow.Year()
	month := int(timeNow.Month())

	switch {
	case month >= 2 && month <= 8: // interval from February to August
		month = 2 // study start on February
		day = 7   // 7.02 - the first day
		endSemester = fmt.Sprintf("%v", year) + endSummerSemester
	case month >= 9 && month <= 12 || month == 1: // interval from September to January
		month = 9 // study start on September
		day = 1   // 1.09 - the first day
		endSemester = fmt.Sprintf("%v", year) + endWinterSemester
	}

	//	calendarId := "primary" // Use account calendar
	calendarId := insertedCalendar.Id

	allWeek, err := dbExplorer(db, group)
	if err != nil {
		return err
	}

	tNow := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	firstDay := int(tNow.Weekday()) - 1
	tNow = tNow.AddDate(0, 0, 7-firstDay)

	if firstDay == -1 {
		firstDay = 0
	}

	for j, allDay := range allWeek {
		if j == firstDay {
			isOdd = true
			tNow = tNow.AddDate(0, 0, -7)
		}

		lessonStart := tNow.Format("2006-01-02")
		sInfo := SubjectsInfo{
			IsOdd:         isOdd,
			LessonStartAt: lessonStart,
			TimeNow:       tNow,
			SemesterEnd:   endSemester,
		}

		var counts int
		for i, lesson := range allDay {
			sInfo.Subject = lesson
			sInfo.Number = i

			events, isEmpty := sInfo.parseAt()
			if isEmpty {
				continue
			}

			for _, event := range events {
				_, err = srvc.Events.Insert(calendarId, event).Do()
				if err != nil {
					return err
				}
				counts++
				// log.Printf("Event created: %s\n", event.HtmlLink)
			}
		}
		log.Printf("%v events created at %v with account %v", counts, tNow.Weekday(), u.Email)

		tNow = tNow.AddDate(0, 0, 1)
	}

	return nil
}

func getSubjects(subj Subject) []Subject {
	if strings.Contains(subj.Name, "#") {
		return subj.parseSharp()
	} else {
		return []Subject{subj}
	}
}

func (sInfo *SubjectsInfo) createEvent() *calendar.Event {
	endSemester := sInfo.SemesterEnd
	freq := make([]string, 0, 1)
	if sInfo.IsAllDay {
		freq = []string{"RRULE:FREQ=WEEKLY;UNTIL=" + endSemester}
	} else {
		freq = []string{"RRULE:FREQ=WEEKLY;INTERVAL=2;UNTIL=" + endSemester}
	}

	lesson := sInfo.Subject
	lesson.Room = getEmpty(lesson.Room)
	lesson.Lector = getEmpty(lesson.Lector)

	if _, isNorth := north[lesson.Room]; isNorth {
		lesson.Room = lesson.Room + "(СЕВЕР)"
	}
	if _, isSouth := south[lesson.Room]; isSouth {
		lesson.Room = lesson.Room + "(ЮГ)"
	}

	i := sInfo.Number
	lessonStart := sInfo.LessonStartAt

	return &calendar.Event{
		Summary:     lesson.Room + " " + lesson.Name + " " + lesson.Lector,
		Location:    "Lomonosov Moscow State University",
		Description: lesson.Lector,
		Start: &calendar.EventDateTime{
			DateTime: lessonStart + timeIntervals[i].Start, // TODO: spring ----> season
			TimeZone: "Europe/Moscow",
		},
		End: &calendar.EventDateTime{
			DateTime: lessonStart + timeIntervals[i].End,
			TimeZone: "Europe/Moscow",
		},
		ColorId: getColorId(lesson.Name, lesson.Room),
		Reminders: &calendar.EventReminders{
			UseDefault: false,
			Overrides:  []*calendar.EventReminder{},
			// ForceSendFields is required, if you dont want to set up notifications, because
			// by default, empty values are omitted from API requests
			ForceSendFields: []string{"UseDefault", "Overrides"},
		},
		Recurrence: freq,
	}
}

func getEmpty(in string) string {
	if in == "__" {
		return ""
	}
	return in
}

func getColorId(name, room string) string {
	/*
		ColorId : Color
		1 : lavender
		2 : sage //шалфей
		3 : grape
		4 : flamingo
		5 : banana
		6 : mandarin
		7 : peacock //павлин
		8 : graphite
		9 : blueberry
		10 : basil //базилик
		11 : tomato
	*/
	switch {
	case name == war:
		return "11"
	case name == practice:
		return "10"
	case name == mfk || name == MFKabbr || name == MFK:
		return "4"
	}

	_, isLecture := audience[room]
	if reUpp.MatchString(name) || isLecture {
		return "3"
	}

	if strings.Contains(name, "с/к") || strings.Contains(name, "НИС") ||
		strings.Contains(name, "ДМП") || strings.Contains(name, "Д/п") ||
		strings.Contains(name, "Д/П") || strings.Contains(name, "C/К") ||
		strings.Contains(name, "С/К") || strings.Contains(name, "ФТД") {
		return "2"
	}

	return "7"
}
