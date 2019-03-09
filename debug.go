/*
	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List("primary").ShowDeleted(false).SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	fmt.Println("Upcoming events:")
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			fmt.Printf("%v (%v)\n", item.Summary, date)
		}
	}
*/

// Refer to the Go quickstart on how to setup the environment:
	// https://developers.google.com/calendar/quickstart/go
	// Change the scope to calendar.CalendarScope and delete any stored credentials.
	event := &calendar.Event{
		Summary:     "Name of Subject and Lecturer",
		Location:    "Lomonosov Moscow State University",//Number of room and direction?
		Description: "Lecturer's name",
		Start: &calendar.EventDateTime{
			DateTime: "2019-03-09T22:22:00+03:00",
			TimeZone: "Europe/Moscow",
		},
		End: &calendar.EventDateTime{
			DateTime: "2019-03-09T23:23:00+03:00",
			TimeZone: "Europe/Moscow",
		},
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
		ColorId: "10",
		Reminders: &calendar.EventReminders{
			UseDefault: false,
			Overrides: []*calendar.EventReminder{},
			//ForceSendFields is required, if you dont want to set up notifications, because
			//by default, empty values are omitted from API requests
			ForceSendFields: []string{"UseDefault", "Overrides"},
		},
		Recurrence: []string{"RRULE:FREQ=WEEKLY;INTERVAL=2;UNTIL=20190601"},
		/*Attendees: []*calendar.EventAttendee{
			&calendar.EventAttendee{Email: "lpage@example.com"},
			&calendar.EventAttendee{Email: "sbrin@example.com"},
		},*/
	}
