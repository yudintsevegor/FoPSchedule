package main

import (
	"strings"

	"fopSchedule/master/common"

	"google.golang.org/api/calendar/v3"
)

// delimiters descriprion:
// @ - odd/even week
// # - lessons in the same time, ex.: English lesson, in which teacher devides groups on 2 subgroups
// % - Name%Lector%Room

func parsePercent(arr []string) []Subject {
	result := make([]Subject, 0, 5)
	for _, val := range arr {
		res := strings.Split(val, "%")
		result = append(result, Subject{
			Name:   res[0],
			Lector: res[1],
			Room:   res[2],
		})
	}

	return result
}

func (sInfo SubjectsInfo) parseAt() ([]*calendar.Event, bool) {
	rawSubjects := sInfo.Subject
	if rawSubjects.Name == "" || rawSubjects.Name == "__" {
		return nil, true
	}

	if strings.Contains(common.LessonCases, rawSubjects.Name) {
		return nil, true
	}

	result := make([]*calendar.Event, 0, 1)
	if !strings.Contains(rawSubjects.Name, "@") {
		if rawSubjects.Name == common.Practice {
			return nil, true
		}

		sInfo.IsAllDay = true
		subjects := getSubjects(rawSubjects)
		for _, subj := range subjects {
			sInfo.Subject = subj
			result = append(result, sInfo.createEvent())
		}

		return result, false
	}

	sInfo.IsAllDay = false

	names := strings.Split(rawSubjects.Name, "@")
	lectors := strings.Split(rawSubjects.Lector, "@")
	rooms := strings.Split(rawSubjects.Room, "@")

	var (
		oddLessonStart  string
		evenLessonStart string
	)

	tNow := sInfo.TimeNow
	lessonStart := tNow.Format(common.TimeLayout)

	if !sInfo.IsOdd {
		oddLessonStart = tNow.AddDate(0, 0, 7).Format(common.TimeLayout)
		evenLessonStart = lessonStart
	} else {
		oddLessonStart = lessonStart
		evenLessonStart = tNow.AddDate(0, 0, 7).Format(common.TimeLayout)
	}

	oddSubject := Subject{
		Name:          names[0],
		Lector:        lectors[0],
		Room:          rooms[0],
		LessonStartAt: oddLessonStart,
	}

	evenSubject := Subject{
		Name:          names[1],
		Lector:        lectors[1],
		Room:          rooms[1],
		LessonStartAt: evenLessonStart,
	}

	oneDay := []Subject{oddSubject, evenSubject}
	for _, subj := range oneDay {
		subjects := getSubjects(subj)
		for _, subj := range subjects {
			if subj.Name != "" && subj.Name != "__" && subj.Name != common.Practice {
				sInfo.LessonStartAt = subj.LessonStartAt
				sInfo.Subject = subj

				result = append(result, sInfo.createEvent())
			}
		}
	}

	return result, false
}

func (subj *Subject) parseSharp() []Subject {
	names := strings.Split(subj.Name, "#")
	lectors := strings.Split(subj.Lector, "#")
	rooms := strings.Split(subj.Room, "#")

	subjects := make([]Subject, 0, len(names))
	for i := 0; i < len(names); i++ {
		subjects = append(subjects, Subject{
			Name:          names[i],
			Room:          rooms[i],
			Lector:        lectors[i],
			LessonStartAt: subj.LessonStartAt,
		})
	}

	return subjects
}
