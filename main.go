package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"regexp"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	///config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	// ====================================================================
	// Get data from database
	db, err := sql.Open("mysql", DSN)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	group := "316б"
	allWeek := dbExplorer(db, group)

	clndr := &calendar.Calendar{
		Summary: "Shedule" + group,
	}
	insertedCalendar, err := srv.Calendars.Insert(clndr).Do()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("==========")
	//	calendarId := "primary" // Use account shedule
	calendarId := insertedCalendar.Id

	var isOdd bool
	var day int
	var endSemester string

	var endSummer = "0601"
	var endWinter = "1231"

	timeNow := time.Now()
	year := timeNow.Year()
	month := int(timeNow.Month())

	if month > 0 && month < 8 {
		month = 2
		day = 7
		endSemester = fmt.Sprintf("%v", year) + endSummer
	} else if month > 7 && month < 13 {
		month = 9
		day = 1
		endSemester = fmt.Sprintf("%v", year) + endWinter
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
			st := DataToParsingAt{
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
					log.Fatalf("Unable to create event. %v\n", err)
				}
				fmt.Printf("Event created: %s\n", event.HtmlLink)
			}
		}
		t = t.AddDate(0, 0, 1)
	}
}

func (st *Subject) parseSharp() ([]Subject){
	count := strings.Count(st.Name, "#")
	str := strings.Repeat("(.*)#", count) + "(.*)"
	reSharp := regexp.MustCompile(str)
	names := reSharp.FindStringSubmatch(st.Name)[1:count+2]
	lectors := reSharp.FindStringSubmatch(st.Lector)[1:count+2]
	rooms := reSharp.FindStringSubmatch(st.Room)[1:count+2]
	
	var subjects = make([]Subject, 0, 1)
	for i := 0; i < len(names); i++{
		subj := Subject{
			Name: names[i],
			Room: rooms[i],
			Lector: lectors[i],
			Parity: st.Parity,
		}
		subjects = append(subjects, subj)
	}
	return subjects
}

func (st *DataToParsingAt) parseAt() ([]*calendar.Event, bool) {
	subject := st.Lesson
	isOdd := st.Parity
	lessonStart := st.StartTime
	t := st.Time

	var result = make([]*calendar.Event, 0, 2)
	if subject.Name == "" || subject.Name == "__" {
		return result, true
	}

	if strings.Contains(cases, subject.Name){
		return result, true
	}

	
	if strings.Contains(subject.Name, "@") {
		st.IsAllDay = false

		regName := reAt.FindStringSubmatch(subject.Name)
		regLector := reAt.FindStringSubmatch(subject.Lector)
		regRoom := reAt.FindStringSubmatch(subject.Room)

		var oddLessonStart string
		var evenLessonStart string

		if isOdd {
			oddLessonStart = lessonStart
			evenLessonStart = t.AddDate(0, 0, 7).Format("2006-01-02")
		} else {
			oddLessonStart = t.AddDate(0, 0, 7).Format("2006-01-02")
			evenLessonStart = lessonStart
		}
		
		oddSubject := Subject{Name: regName[1], Lector: regLector[1], Room: regRoom[1], Parity: oddLessonStart}
		evenSubject := Subject{Name: regName[2], Lector: regLector[2], Room: regRoom[2], Parity: evenLessonStart}

		var arr = []Subject{oddSubject, evenSubject}
		for _, subj := range arr {
			var fromSharp = make([]Subject, 0, 1)
			if strings.Contains(subj.Name, "#"){
				fromSharp = subj.parseSharp()
			} else {
				fromSharp = append(fromSharp, subj)
			}
			for _, v := range fromSharp{
				if v.Name != "" && v.Name != "__" && v.Name != practice {
					st.StartTime = v.Parity
					st.Lesson = v
					event := st.createEvent()
					result = append(result, event)
				}
			}
		}
		return result, false
	}

	var fromSharp = make([]Subject, 0, 1)
	if subject.Name == practice {
		return result, true
	}
	if strings.Contains(subject.Name, "#"){
		fromSharp = subject.parseSharp()
	} else {
		fromSharp = append(fromSharp, subject)
	}
	st.IsAllDay = true
	for _, v := range fromSharp{
		st.Lesson = v
		event := st.createEvent()
		result = append(result, event)
	}
	
	return result, false
}

func (st *DataToParsingAt) createEvent() *calendar.Event {
	subject := st.Lesson
	i := st.Number
	allDay := st.IsAllDay
	lessonStart := st.StartTime
	endSemester := st.SemesterEnd

	var freq = make([]string, 0, 1)
	if allDay {
		freq = []string{"RRULE:FREQ=WEEKLY;UNTIL=" + endSemester}
	} else {
		freq = []string{"RRULE:FREQ=WEEKLY;INTERVAL=2;UNTIL=" + endSemester}
	}
	color := getColorId(subject.Name, subject.Room)
	if subject.Lector == "__" {
		subject.Lector = ""
	}
	if subject.Room == "__" {
		subject.Room = ""
	}
	if _, isNorth := north[subject.Room]; isNorth{
		subject.Room = subject.Room + "(СЕВЕР)"
	}
	if _, isSouth := south[subject.Room]; isSouth{
		subject.Room = subject.Room + "(ЮГ)"
	}
	event := &calendar.Event{
		Summary:     subject.Room + " " + subject.Name + " " + subject.Lector,
		Location:    "Lomonosov Moscow State University", //Number of room and direction?
		Description: subject.Lector,
		Start: &calendar.EventDateTime{
			DateTime: lessonStart + timeIntervals[i].Start, // spring ----> season
			TimeZone: "Europe/Moscow",
		},
		End: &calendar.EventDateTime{
			DateTime: lessonStart + timeIntervals[i].End,
			TimeZone: "Europe/Moscow",
		},
		ColorId: color,
		Reminders: &calendar.EventReminders{
			UseDefault: false,
			Overrides:  []*calendar.EventReminder{},
			//ForceSendFields is required, if you dont want to set up notifications, because
			//by default, empty values are omitted from API requests
			ForceSendFields: []string{"UseDefault", "Overrides"},
		},
		Recurrence: freq,
	}
	return event
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
	if name == war {
		return "11"
	} else if name == practice {
		return "10"
	} else if name == mfk || name == MFKabbr || name == MFK {
		return "4"
	}
	_, isLecture := audience[room]
	if reUpp.MatchString(name) || isLecture {
		return "3"
	}
	if strings.Contains(name, "с/к") || strings.Contains(name, "НИС") || strings.Contains(name, "ДМП") || strings.Contains(name, "Д/п") || strings.Contains(name, "Д/П") || strings.Contains(name, "C/К") || strings.Contains(name, "С/К") || strings.Contains(name, "ФТД") {
		return "2"
	}
	return "7"
}
