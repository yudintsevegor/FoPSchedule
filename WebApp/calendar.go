package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"strings"
	"time"

	"google.golang.org/api/calendar/v3"
)

func putData(client *http.Client, group string) error {
	db, err := sql.Open("mysql", DSN)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	clndr := &calendar.Calendar{
		Summary: calendarName + group,
	}

	srv, err := calendar.New(client)
	if err != nil {
		return err
	}

	insertedCalendar, err := srv.Calendars.Insert(clndr).Do()
	if err != nil {
		return err
	}

	var (
		isOdd       bool
		day         int
		endSemester string
	)

	const (
		endSummer = "0601"
		endWinter = "1231"
	)

	timeNow := time.Now()
	year := timeNow.Year()
	month := int(timeNow.Month())

	switch {
	case month > 0 && month < 8:
		month = 2
		day = 7
		endSemester = fmt.Sprintf("%v", year) + endSummer
	case month > 7 && month < 13:
		month = 9
		day = 1
		endSemester = fmt.Sprintf("%v", year) + endWinter
	}
	/*
		if  month > 0 && month < 8{
			month = 2
			day = 7
			endSemester = fmt.Sprintf("%subj", year) + endSummer
		} else if month > 7 && month < 13 {
			month = 9
			day = 1
			endSemester = fmt.Sprintf("%subj", year) + endWinter
		}
	*/

	//	calendarId := "primary" // Use account calendar
	calendarId := insertedCalendar.Id

	allWeek, err := dbExplorer(db, group)
	if err != nil {
		return err
	}
	
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	firstDay := int(t.Weekday()) - 1
	t = t.AddDate(0, 0, 7-firstDay)

	for j := 0; j < len(allWeek); j++ {
		if j == firstDay {
			isOdd = !isOdd
			t = t.AddDate(0, 0, -7)
		}

		lessonStart := t.Format("2006-01-02")
		for i, lesson := range allWeek[j] {
			st := SubjectsInfo{
				Lesson:      lesson,
				Number:      i,
				Parity:      isOdd,
				StartTime:   lessonStart,
				Time:        t,
				SemesterEnd: endSemester,
			}

			events, isEmpty := st.parseAt()
			if isEmpty {
				continue
			}

			for _, event := range events {
				event, err = srv.Events.Insert(calendarId, event).Do()
				if err != nil {
					return err
				}
				log.Printf("Event created: %s\n", event.HtmlLink)
			}
		}

		t = t.AddDate(0, 0, 1)
	}

	return nil
}

func (st *SubjectsInfo) parseAt() ([]*calendar.Event, bool) {
	rawSubjects := st.Lesson
	isOdd := st.Parity
	lessonStart := st.StartTime
	t := st.Time

	result := make([]*calendar.Event, 0, 2)
	if rawSubjects.Name == "" || rawSubjects.Name == "__" {
		return nil, true
	}

	if strings.Contains(lessonCases, rawSubjects.Name) {
		return nil, true
	}

	if !strings.Contains(rawSubjects.Name, "@") {
		if rawSubjects.Name == practice {
			return nil, true
		}

		subjects := make([]Subject, 0, 1)
		if strings.Contains(rawSubjects.Name, "#") {
			subjects = rawSubjects.parseSharp()
		} else {
			subjects = append(subjects, rawSubjects)
		}

		st.IsAllDay = true
		for _, subj := range subjects {
			st.Lesson = subj
			result = append(result, st.createEvent())
		}

		return result, false
	}

	st.IsAllDay = false

	/*
		regName := reAt.FindStringSubmatch(rawSubjects.Name)
		regLector := reAt.FindStringSubmatch(rawSubjects.Lector)
		regRoom := reAt.FindStringSubmatch(rawSubjects.Room)
	*/
	names := strings.Split(rawSubjects.Name, "@")
	lectors := strings.Split(rawSubjects.Lector, "@")
	rooms := strings.Split(rawSubjects.Room, "@")

	var (
		oddLessonStart  string
		evenLessonStart string
	)

	if isOdd {
		oddLessonStart = lessonStart
		evenLessonStart = t.AddDate(0, 0, 7).Format("2006-01-02")
	} else {
		oddLessonStart = t.AddDate(0, 0, 7).Format("2006-01-02")
		evenLessonStart = lessonStart
	}

	oddSubject := Subject{Name: names[0], Lector: lectors[0], Room: rooms[0], Parity: oddLessonStart}
	evenSubject := Subject{Name: names[1], Lector: lectors[1], Room: rooms[1], Parity: evenLessonStart}

	oneDay := []Subject{oddSubject, evenSubject}
	for _, subj := range oneDay {
		subjects := make([]Subject, 0, 1)
		if strings.Contains(subj.Name, "#") {
			subjects = subj.parseSharp()
		} else {
			subjects = append(subjects, subj)
		}

		for _, subj := range subjects {
			if subj.Name != "" && subj.Name != "__" && subj.Name != practice {
				st.StartTime = subj.Parity
				st.Lesson = subj

				result = append(result, st.createEvent())
			}
		}
	}

	return result, false
}

func (st *Subject) parseSharp() []Subject {
	/*
		count := strings.Count(st.Name, "#")
		str := strings.Repeat("(.*)#", count) + "(.*)"

		reSharp := regexp.MustCompile(str)


		names := reSharp.FindStringSubmatch(st.Name)[1 : count+2]
		lectors := reSharp.FindStringSubmatch(st.Lector)[1 : count+2]
		rooms := reSharp.FindStringSubmatch(st.Room)[1 : count+2]
	*/

	names := strings.Split(st.Name, "#")
	lectors := strings.Split(st.Lector, "#")
	rooms := strings.Split(st.Room, "#")

	subjects := make([]Subject, 0, len(names))
	for i := 0; i < len(names); i++ {
		subjects = append(subjects, Subject{
			Name:   names[i],
			Room:   rooms[i],
			Lector: lectors[i],
			Parity: st.Parity,
		})
	}

	return subjects
}

func (st *SubjectsInfo) createEvent() *calendar.Event {
	rawSubjects := st.Lesson
	i := st.Number
	lessonStart := st.StartTime
	endSemester := st.SemesterEnd

	freq := make([]string, 0, 1)
	if st.IsAllDay {
		freq = []string{"RRULE:FREQ=WEEKLY;UNTIL=" + endSemester}
	} else {
		freq = []string{"RRULE:FREQ=WEEKLY;INTERVAL=2;UNTIL=" + endSemester}
	}

	/*
		if rawSubjects.Lector == "__" {
			rawSubjects.Lector = ""
		}
		if rawSubjects.Room == "__" {
			rawSubjects.Room = ""
		}
	*/

	rawSubjects.Room = getEmpty(rawSubjects.Room)
	rawSubjects.Lector = getEmpty(rawSubjects.Lector)

	if _, isNorth := north[rawSubjects.Room]; isNorth {
		rawSubjects.Room = rawSubjects.Room + "(СЕВЕР)"
	}
	if _, isSouth := south[rawSubjects.Room]; isSouth {
		rawSubjects.Room = rawSubjects.Room + "(ЮГ)"
	}

	return &calendar.Event{
		Summary:     rawSubjects.Room + " " + rawSubjects.Name + " " + rawSubjects.Lector,
		Location:    "Lomonosov Moscow State University", //Number of room and direction?
		Description: rawSubjects.Lector,
		Start: &calendar.EventDateTime{
			DateTime: lessonStart + timeIntervals[i].Start, // spring ----> season
			TimeZone: "Europe/Moscow",
		},
		End: &calendar.EventDateTime{
			DateTime: lessonStart + timeIntervals[i].End,
			TimeZone: "Europe/Moscow",
		},
		ColorId: getColorId(rawSubjects.Name, rawSubjects.Room),
		Reminders: &calendar.EventReminders{
			UseDefault: false,
			Overrides:  []*calendar.EventReminder{},
			//ForceSendFields is required, if you dont want to set up notifications, because
			//by default, empty values are omitted from API requests
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
	/*
		if name == war {
			return "11"
		} else if name == practice {
			return "10"
		} else if name == mfk || name == MFKabbr || name == MFK {
			return "4"
		}
	*/

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
