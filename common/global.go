package common

import "regexp"

const (
	MaxConnections = 100

	CookieURL  = "/cookie"
	CookieName = "fopschedule"

	MainHTMLPage  = "WebApp/html/mainPage.html"
	IndexHTMLPage = "WebApp/html/index.html"
	HtmlPath      = "../WebApp/html/index.html"
	HtmlGen       = "newInd.html"

	CalendarName = "Shedule"
	UrlCalendar  = "https://calendar.google.com"
	TimeLayout   = "2006-01-02"

	Columns = " ( first, second, third, fourth, fifth ) "
	QuesStr = " ( ?, ?, ?, ?, ? ) "

	DiplomaPractice    = "Преддипломная практика"
	WarCaps            = "ВОЕННАЯ ПОДГОТОВКА"
	War                = "Военная подготовка"
	MFKCaps            = "МЕЖФАКУЛЬТЕТСКИЕ КУРСЫ"
	MFKabbr            = "МФК"
	Mfk                = "Межфакультетские курсы"
	CommonPhysicsPrac  = "Общий физический практикум"
	SpecPrac           = "Специальный физический практикум"
	SpecPracCaps       = "СПЕЦПРАКТИКУМ"
	RadioPrac          = "Практикум по радиоэлектронике"
	RadioPhysPrac      = "Радиофизика практикум"
	PhysEdu            = "Физическая культура"
	Research           = "Научно-исследовательская практика"
	AstroProblems      = "Современные проблемы астрономии"
	NISShort           = "НИС"
	NIS                = "Научно-исследовательский семинар"
	AstrShort          = "астр."
	IntroInExp         = "Введение в технику эксперимента"
	NuclearPrac        = "Ядерный практикум"
	AtomicPrac         = "Атомный практикум"
	TeacherPractice    = "Педагогическая практика"
	NIW                = "Научно-исследовательская работа"
	ModernphysProblems = "Современные проблемы физики"

	LessonCases = WarCaps + " " + War + " " + MFKCaps + " " + Mfk + " " + MFKabbr
)

var LessonMap = map[string]struct{}{
	DiplomaPractice:    struct{}{},
	WarCaps:            struct{}{},
	War:                struct{}{},
	MFKCaps:            struct{}{},
	MFKabbr:            struct{}{},
	Mfk:                struct{}{},
	CommonPhysicsPrac:  struct{}{},
	SpecPrac:           struct{}{},
	SpecPracCaps:       struct{}{},
	RadioPrac:          struct{}{},
	RadioPhysPrac:      struct{}{},
	PhysEdu:            struct{}{},
	Research:           struct{}{},
	AstroProblems:      struct{}{},
	NISShort:           struct{}{},
	NIS:                struct{}{},
	AstrShort:          struct{}{},
	IntroInExp:         struct{}{},
	NuclearPrac:        struct{}{},
	AtomicPrac:         struct{}{},
	TeacherPractice:    struct{}{},
	NIW:                struct{}{},
	ModernphysProblems: struct{}{},
}

var (
	// to set another color for event if event is for all groups
	ReUpp  = regexp.MustCompile("([А-Я]){5,}")
	ReNum  = regexp.MustCompile(`([0-9]+М*Б*)`)
	ReDash = regexp.MustCompile(`(\s\-\s)`)

	SubGroups = map[string][]string{
		"341":  []string{"341а", "341б"},
		"441":  []string{"441а", "441б"},
		"141М": []string{"141Ма", "141Мб"},
		"241М": []string{"241Ма", "241Мб"},
		"316":  []string{"316а", "316б"},
		"416":  []string{"416а", "416б"},
		"116М": []string{"116Ма", "116Мб"},
		"216М": []string{"216Ма", "216Мб"},
	}

	South = map[string]string{
		"5-23": "",
		"5-24": "",
		"5-25": "",
		"5-26": "",
		"5-27": "",
		"5-38": "",
		"5-39": "",
		"5-40": "",
		"5-41": "",
		"5-42": "",
		"5-18": "",
		"5-19": "",
	}
	North = map[string]string{
		"5-33":   "",
		"5-34":   "",
		"5-35":   "",
		"5-36":   "",
		"5-37":   "",
		"5-44":   "",
		"5-45":   "",
		"5-46":   "",
		"5-47":   "",
		"5-48":   "",
		"5-49":   "",
		"5-50":   "",
		"5-51":   "",
		"5-52":   "",
		"5-53":   "",
		"5-61":   "",
		"5-62":   "",
		"5-68":   "",
		"Л.каб.": "",
	}

	Audience = map[string]string{
		"СФА":              "",
		"ЮФА":              "",
		"ЦФА":              "",
		"Ауд. им. Хохлова": "",
	}

	moscowTime    = "+03:00"
	TimeIntervals = map[int]LessonRange{
		0: {Start: "T9:00:00" + moscowTime, End: "T10:35:00" + moscowTime},
		1: {Start: "T10:50:00" + moscowTime, End: "T12:25:00" + moscowTime},
		2: {Start: "T13:30:00" + moscowTime, End: "T15:05:00" + moscowTime},
		3: {Start: "T15:20:00" + moscowTime, End: "T16:55:00" + moscowTime},
		4: {Start: "T17:05:00" + moscowTime, End: "T18:40:00" + moscowTime},
	}
)

type LessonRange struct {
	Start string
	End   string
}
