package main

type Subject struct {
	Name   string
	Lector string
	Room   string
}

type Department struct {
	Number  string
	Lessons []Subject
}

type Interval struct {
	Start int
	End   int
}

type LessonRange struct {
	Start string
	End string
}
