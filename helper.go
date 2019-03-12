package main

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func (st *DataToParsingLine) parseLine(subjectIndex, countSmall0 int, text string, nextLine, is2Weeks, isFirstInSmall0 bool) ([]Department, []string) {
	departments := st.Departments
	allGr := st.AllGroups
	resFromReg := st.ResultFromReqexp
	insertedGroups := st.InsertedGroups
	subject := st.Lesson
	reInterval := st.RegexpInterval

//	fmt.Println(allGr)
	if len(resFromReg) == 0 {
		for _, dep := range departments {
			for _, gr := range allGr {
				if dep.Number != gr {
					continue
				}
				if !nextLine {
					if countSmall0 < 0 {
						dep.Lessons[subjectIndex] = subject
						continue
					}
					newSubj := Subject{}
					if isFirstInSmall0 {
						newSubj = subject
					} else {
						newSubj = Subject{
							Name: dep.Lessons[subjectIndex].Name + "#" + subject.Name,
							Lector: dep.Lessons[subjectIndex].Lector + "#" + subject.Lector,
							Room: dep.Lessons[subjectIndex].Room + "#" + subject.Room,
						}
					}
					dep.Lessons[subjectIndex] = newSubj
					continue
				}
				newSubj := subject.getNewStruct(dep.Lessons[subjectIndex])
				dep.Lessons[subjectIndex] = newSubj
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
			if !nextLine {
				if dep.Lessons[subjectIndex].Name == ""{
					dep.Lessons[subjectIndex] = subject
				} else {
					newSubj := Subject{
						Name: dep.Lessons[subjectIndex].Name + "#" + subject.Name,
						Lector: dep.Lessons[subjectIndex].Lector + "#" + subject.Lector,
						Room: dep.Lessons[subjectIndex].Room + "#" + subject.Room,
					}
					dep.Lessons[subjectIndex] = newSubj
				}
				insertedGroups = append(insertedGroups, gr)
				continue
			}
			newSubj := subject.getNewStruct(dep.Lessons[subjectIndex])
			dep.Lessons[subjectIndex] = newSubj
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

	for _, dep := range departments {
		for gr, _ := range mapAllGr {
			if dep.Number != gr {
				continue
			}
			newSubj := Subject{
				Name:   "__",
				Lector: "__",
				Room:   "__",
			}
			if nextLine {
				newSubj = newSubj.getNewStruct(dep.Lessons[subjectIndex])
			}
			dep.Lessons[subjectIndex] = newSubj
		}
	}
	if countSmall0 <= 0 {
		insertedGroups = make([]string, 5)
	}

	return departments, insertedGroups
}

func putToDB(departments []Department, db *sql.DB) {
	for _, val := range departments {
		var valuesToDB = make([]interface{}, 0, 1)
		for _, les := range val.Lessons {
			value := les.Name + "%" + les.Lector + "%" + les.Room
			valuesToDB = append(valuesToDB, value)
		}
		//		fmt.Println(valuesToDB)
		req := fmt.Sprintf("INSERT INTO `%v`"+columns+"VALUES"+quesStr, val.Number)
		statement, err := db.Prepare(req)
		if err != nil {
			log.Fatal(err)
		}
		_, err = statement.Exec(valuesToDB...)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(val.Number)
		fmt.Println("PUT IN TABLE")
	}
}

func clean(arr []Department) []Department {
	var result = make([]Department, len(arr))
	for i, d := range arr {
		s := Department{}
		s.Number = d.Number
		s.Lessons = make([]Subject, len(d.Lessons))
		result[i] = s
	}

	return result
}

func fromStringToInt(class string) int {
	num := reNum.FindStringSubmatch(class)[1]
	numberFromClass, err := strconv.Atoi(num)
	if err != nil {
		log.Fatal(err)
	}

	return numberFromClass
}

func (st *Subject) getNewStruct(subject Subject) Subject {
	return Subject{
		Name:   subject.Name + "@" + st.Name,
		Lector: subject.Lector + "@" + st.Lector,
		Room:   subject.Room + "@" + st.Room,
	}
}

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
	if strings.Contains(text, MFK) {
		subj.Name = MFK
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
	if strings.Contains(text, prac201) {
		subj.Name = text
		subj.Room = "__"
		subj.Lector = "__"
		return subj
	}
	if strings.Contains(text, prac){
		subj.Name = prac
		subj.Room = "__"
		subj.Lector = "__"
		return subj
	}
	if strings.Contains(text, specprac) {
		subj.Name = text 
		subj.Room = "__"
		subj.Lector = "__"
		return subj
	}
	if strings.Contains(text, phys){
		subj.Name = phys
		subj.Room = "__"
		subj.Lector = "__"
		return subj
	}
	fmt.Println(text)
	rLect := regexp.MustCompile(`.* ` + room + ` (.*)`)
	Lect := rLect.FindStringSubmatch(text)[1]

	rSubj := regexp.MustCompile(`([^0-9\-]*) ` + room + " " + Lect)
	Subj := rSubj.FindStringSubmatch(text)[1]

	subj.Name = Subj
	subj.Lector = Lect
	subj.Room = room

	return subj
}
