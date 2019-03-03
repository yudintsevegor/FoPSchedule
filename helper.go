package main

import(
	"regexp"
	"strings"
	"unicode"
	"strconv"
)

var course = "4"
var reInterval = regexp.MustCompile(`(` +  course + `\d{2})\s*\-\s*` + `(` +  course + `\d{2})`)

func parseLine(departments []Department,  allGr []string, resFromReg []string, subject Subject, text string, n int, nextStr bool) ([]Department){
	
	if len(resFromReg) == 0{
		for _, dep := range departments{
			for _, gr := range allGr{
				if dep.Number != gr {
					continue
				}
				if !nextStr {
					dep.Lessons[n] = subject
				} else {
					newSubj := Subject{}
					newSubj.Lector = dep.Lessons[n].Lector + "@" + subject.Lector
					newSubj.Room = dep.Lessons[n].Room + "@" + subject.Room
					newSubj.Name = dep.Lessons[n].Name + "@" + subject.Name
					dep.Lessons[n] = newSubj
				}
			}
		}
	} else {
		if reInterval.MatchString(text){
			interval := reInterval.FindStringSubmatch(text)
			left, _ := strconv.Atoi(interval[1])
			right, _ := strconv.Atoi(interval[2])
			
			for i := left + 1; i < right; i++{
				resFromReg = append(resFromReg, strconv.Itoa(i))
			}
		}
		for _, dep := range departments{
			for _, gr := range resFromReg{
				if dep.Number != gr {
					continue
				}
				if !nextStr{
					dep.Lessons[n] = subject
				} else {
					newSubj := Subject{}
					newSubj.Lector = dep.Lessons[n].Lector + "@" + subject.Lector
					newSubj.Room = dep.Lessons[n].Room + "@" + subject.Room
					newSubj.Name = dep.Lessons[n].Name + "@" + subject.Name
					dep.Lessons[n] = newSubj
				}
			}
		}
	}
	
	return departments
}

func fromStringToInt(class string) (int, error) {
	num := re.FindStringSubmatch(class)[1]
	number, err := strconv.Atoi(num)

	return number, err
}

var practice = "Преддипломная практика"
var war = "ВОЕННАЯ ПОДГОТОВКА"
var mfk = "МЕЖФАКУЛЬТЕТСКИЕ КУРСЫ"

func parseGroups(text, room string) Subject {
	subj := Subject{}

	var isSpace bool
	for _, val := range text {
		isSpace = unicode.IsSpace(val)
		break
	}

	if isSpace {
		subj.Name = "__"
		return subj
	}

	if strings.Contains(text, practice) {
		subj.Name = practice

		return subj
	}

	if strings.Contains(text, mfk) {
		subj.Name = mfk

		return subj
	}
	if strings.Contains(text, war) {
		subj.Name = war

		return subj
	}

	//	fmt.Println("TEXT: ", text)
	rLect := regexp.MustCompile(`.* ` + room + ` (.*)`)
	Lect := rLect.FindStringSubmatch(text)[1]

	rSubj := regexp.MustCompile(`([^0-9\-]*) ` + room + " " + Lect)
	Subj := rSubj.FindStringSubmatch(text)[1]

	subj.Name = Subj
	subj.Lector = Lect
	subj.Room = room


	return subj
}
