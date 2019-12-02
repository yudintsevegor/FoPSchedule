package common

import "strings"

func (subj *Subject) ParseSharp() []Subject {
	names := strings.Split(subj.Name, "#")
	lectors := strings.Split(subj.Lector, "#")
	rooms := strings.Split(subj.Room, "#")

	var subjects = make([]Subject, 0, 1)
	for i := 0; i < len(names); i++ {
		subj := Subject{
			Name:          names[i],
			Room:          rooms[i],
			Lector:        lectors[i],
			Parity:        subj.Parity,
			LessonStartAt: subj.LessonStartAt,
		}
		subjects = append(subjects, subj)
	}

	return subjects
}
