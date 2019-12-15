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

func (st *Subject) GetNewStruct(subject Subject, delimiter string) Subject {
	return Subject{
		Name:   subject.Name + delimiter + st.Name,
		Lector: subject.Lector + delimiter + st.Lector,
		Room:   subject.Room + delimiter + st.Room,
	}
}
