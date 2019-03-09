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

	// Get data from database
	db, err := sql.Open("mysql", DSN)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	group := "442"
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
	year := 2019
	var month time.Month = 2
	day := 7
	t := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	firstDay := int(t.Weekday()) - 1
	t = t.AddDate(0, 0, 7-firstDay)

	for j := 0; j < len(allWeek); j++ {
		if j == firstDay {
			isOdd = !isOdd
			t = t.AddDate(0, 0, -7)
		}
		lessonStart := t.Format("2006-01-02")
		for i, lesson := range allWeek[j] {
			events, isEmpty := parseAt(lesson, i, isOdd, lessonStart, t)
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

var springStart = "2019-02-07"
var springEnd = "20190601"

//var autumnEnd = "20191231"

func createEvent(subject Subject, i int, allDay bool, lessonStart string) *calendar.Event {
	var freq = make([]string, 0, 1)
	if allDay {
		freq = []string{"RRULE:FREQ=WEEKLY;UNTIL=" + springEnd}
	} else {
		freq = []string{"RRULE:FREQ=WEEKLY;INTERVAL=2;UNTIL=" + springEnd}
	}
	color := getColorId(subject.Name)
	if subject.Lector == "__" {
		subject.Lector = ""
	}
	if subject.Room == "__" {
		subject.Room = ""
	}

	event := &calendar.Event{
		Summary:     subject.Name + " " + subject.Lector + " " + subject.Room,
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

func parseAt(subject Subject, i int, isOdd bool, lessonStart string, t time.Time) ([]*calendar.Event, bool) {
	var result = make([]*calendar.Event, 0, 2)
	if subject.Name == "" || subject.Name == "__" {
		return result, true
	}

	if subject.Name == war || subject.Name == mfk {
		return result, true
	}

	if strings.Contains(subject.Name, "@") {
		regName := reAt.FindStringSubmatch(subject.Name)
		regLector := reAt.FindStringSubmatch(subject.Lector)
		regRoom := reAt.FindStringSubmatch(subject.Room)

		oddSubject := Subject{Name: regName[1], Lector: regLector[1], Room: regRoom[1]}
		evenSubject := Subject{Name: regName[2], Lector: regLector[2], Room: regRoom[2]}

		var oddLessonStart string
		var evenLessonStart string

		if isOdd {
			oddLessonStart = lessonStart
			evenLessonStart = t.AddDate(0, 0, 7).Format("2006-01-02")
		} else {
			oddLessonStart = t.AddDate(0, 0, 7).Format("2006-01-02")
			evenLessonStart = lessonStart
		}

		if oddSubject.Name != "" && oddSubject.Name != "__" && oddSubject.Name != practice {
			event := createEvent(oddSubject, i, false, oddLessonStart)
			result = append(result, event)
		}
		if evenSubject.Name != "" && evenSubject.Name != "__" && evenSubject.Name != practice {
			event := createEvent(evenSubject, i, false, evenLessonStart)
			result = append(result, event)
		}
		return result, false
	}

	if subject.Name == practice {
		return result, true
	}
	event := createEvent(subject, i, true, lessonStart)
	result = append(result, event)

	return result, false
}

func getColorId(name string) string {
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
	} else if name == mfk {
		return "4"
	}
	if reUpp.MatchString(name) {
		return "3"
	}
	if strings.Contains(name, "Д/П") || strings.Contains(name, "С/К") {
		return "2"
	}
	return "7"
}
