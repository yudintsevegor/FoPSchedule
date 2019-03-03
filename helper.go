package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var course = "4"
var reInterval = regexp.MustCompile(`(` + course + `\d{2})\s*\-\s*` + `(` + course + `\d{2})`)

func parseLine(departments []Department, allGr []string, resFromReg []string, insertedGroups []string, subject Subject, text string, countSmall0, n int, nextStr, is2Weeks bool) ([]Department, []string) {
	if len(resFromReg) == 0 {
		for _, dep := range departments {
			for _, gr := range allGr {
				if dep.Number != gr {
					continue
				}
				if !nextStr {
					dep.Lessons[n] = subject
					continue
				}
				newSubj := Subject{}
				newSubj.Lector = dep.Lessons[n].Lector + "@" + subject.Lector
				newSubj.Room = dep.Lessons[n].Room + "@" + subject.Room
				newSubj.Name = dep.Lessons[n].Name + "@" + subject.Name
				dep.Lessons[n] = newSubj
			}
		}
		return departments, insertedGroups
	}

	if reInterval.MatchString(text) {
		interval := reInterval.FindStringSubmatch(text)
		left, _ := strconv.Atoi(interval[1])
		right, _ := strconv.Atoi(interval[2])

		for i := left + 1; i < right; i++ {
			resFromReg = append(resFromReg, strconv.Itoa(i))
		}
	}
	for _, dep := range departments {
		for _, gr := range resFromReg {
			if dep.Number != gr {
				continue
			}
			if !nextStr {
				dep.Lessons[n] = subject
				insertedGroups = append(insertedGroups, gr)
				continue
			}
			newSubj := Subject{}
			newSubj.Lector = dep.Lessons[n].Lector + "@" + subject.Lector
			newSubj.Room = dep.Lessons[n].Room + "@" + subject.Room
			newSubj.Name = dep.Lessons[n].Name + "@" + subject.Name
			dep.Lessons[n] = newSubj
			insertedGroups = append(insertedGroups, gr)
		}
	}

	if countSmall0 > 0 || is2Weeks {
		return departments, insertedGroups
	}

	var mapAllGr = make(map[string]string)
	for _, gr := range allGr {
		mapAllGr[gr] = ""
	}

	for _, v1 := range insertedGroups {
		for v2, _ := range mapAllGr {
			if v2 == v1 {
				delete(mapAllGr, v2)
			}
		}
	}
	//	fmt.Println("=================inserted groups and mapAllGr=========================================")
	//	fmt.Println(insertedGroups)
	//	fmt.Println("==========================================================")
	//	fmt.Println(mapAllGr)
	//	fmt.Println("==========================================================")

	for _, dep := range departments {
		for gr, _ := range mapAllGr {
			if dep.Number != gr {
				continue
			}
			newSubj := Subject{}
			if !nextStr {
				newSubj.Name = "__"
				newSubj.Room = "__"
				newSubj.Lector = "__"
			} else {
				newSubj.Name = dep.Lessons[n].Name + "@" + "__"
				newSubj.Room = dep.Lessons[n].Room + "@" + "__"
				newSubj.Lector = dep.Lessons[n].Lector + "@" + "__"
			}

			dep.Lessons[n] = newSubj
		}
	}
	if countSmall0 <= 0 {
		insertedGroups = make([]string, 5)
	}

	return departments, insertedGroups
}

func fromStringToInt(class string) int {
	num := re.FindStringSubmatch(class)[1]
	number, err := strconv.Atoi(num)
	if err != nil {
		log.Fatal(err)
	}

	return number
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
		subj.Room = "__"
		subj.Lector = "__"
		return subj
	}

	if strings.Contains(text, practice) {
		subj.Name = practice
		subj.Room = "__"
		subj.Lector = "__"
		return subj
	}

	if strings.Contains(text, mfk) {
		subj.Name = mfk
		subj.Room = "__"
		subj.Lector = "__"
		return subj
	}
	if strings.Contains(text, war) {
		subj.Name = war
		subj.Room = "__"
		subj.Lector = "__"
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
