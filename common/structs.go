package common

type LessonRange struct {
	Start string
	End   string
}

type Subject struct {
	Name   string
	Lector string
	Room   string

	Parity        string
	LessonStartAt string
}
