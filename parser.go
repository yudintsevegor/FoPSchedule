package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var courses = map[string][]string{
//		"4": []string{"2"},
		"1": []string{"1", "2", "3"},
		"2": []string{"1", "2", "3"},
		"3": []string{"1", "2"},
		"4": []string{"1", "2"},
		"5": []string{"1", "2"},
		"6": []string{"1", "2"},
	}
	db, err := sql.Open("mysql", DSN)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	for course, arr := range courses{
		for _, thread := range arr {
			res, err := http.Get("http://ras.phys.msu.ru/table/" + course + "/" + thread + ".html")
			if err != nil {
				log.Fatal(err)
			}
			defer res.Body.Close()
			if res.StatusCode != 200 {
				log.Fatal("status code error: %d %s", res.StatusCode, res.Status)
			}
			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				log.Fatal(err)
			}
			parse(course, db, doc)
		}
	}
}

func parse(course string, db *sql.DB, doc *goquery.Document) {
	if course == "5"{
		course = "[1,5]"
	}
	if course == "6"{
		course = "[2,6]"
	}
	
	var reGrp = regexp.MustCompile(course + `\d{2}[МБаб]*`)
	var reInterval = regexp.MustCompile(`(` + course + `\d{2}[МБаб]*)\s*\-\s*` + `(` + course + `\d{2}[МБаб]*)`)

	grpBegin := "ГРУППЫ >>"
	grpEnd := "<< ГРУППЫ"

	var err error
	var grpsFound int
	var isGroups bool
	var departments = make([]Department, 0, 5)
	var eachColumn = make(map[int][]string)

	columnIndex := 0
	doc.Find("td").Each(func(i int, std *goquery.Selection) {
		if grpsFound > 1 {
			return
		}
		text := std.Text()
		if isGroups && text != grpEnd {
			tmpSlice := reGrp.FindAllString(text, -1)
			resFromReg := make([]string, 0, len(tmpSlice))
			for _, gr := range tmpSlice{
				if subgr, ok := subGroups[gr]; ok{
					resFromReg = append(resFromReg, subgr...)
					continue
				}
				resFromReg = append(resFromReg, gr)
			}
			eachColumn[columnIndex] = resFromReg
			columnIndex++
			for _, val := range resFromReg {
				depart := Department{Lessons: make([]Subject, 5, 5)}
				depart.Number = val
				departments = append(departments, depart)
			}
		}
		if text == grpBegin {
			grpsFound++
			isGroups = true
		} else if text == grpEnd {
			isGroups = false
		}
	})

	for key, val := range eachColumn {
		fmt.Println(key, val)
	}
	partOfReq := `(
				  id int(11) NOT NULL AUTO_INCREMENT,
				  first text(255),
				  second text(255),
				  third text(255),
				  fourth text(255),
				  fifth text(255),
				  PRIMARY KEY (id)
				) ENGINE=InnoDB DEFAULT CHARSET=utf8; `

	for _, val := range eachColumn {
		for _, gr := range val {
			del := fmt.Sprintf("DROP TABLE IF EXISTS `%v`; ", gr)
			_, err = db.Exec(del)
			if err != nil {
				log.Fatal(err)
			}
			request := fmt.Sprintf("CREATE TABLE `%v` "+partOfReq, gr)
			_, err = db.Exec(request)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	fmt.Println("TABLES CREATED")

	var time string
	tditem := "tditem"
	tdsmall := "tdsmall"
	tdtime := "tdtime"

	t := "9:00- - -  10:35"
	var tmp int

	var classBeforeSmall0 string
	var numberBeforeSmall0 int
	var countSmall0 int

	var ind int
	var subjectIndex int
	var spanIndex int

	var nextLine bool
	var is2Weeks bool
	var isFirstInSmall0 bool

	var Spans = make([]Interval, 10, 10)
	var insertedGroups = make([]string, 5)

	var Saturday int
	isSaturday := false

	doc.Find("td").Each(func(i int, std *goquery.Selection) {
		text := std.Text()

		if text == grpBegin {
			// there are 3 <tr> with groups
			Saturday++
		}

		if Saturday == 3 && !isSaturday {
			putToDB(departments, db)
			isSaturday = true
		}

		class, ok := std.Attr("class")
		if !ok {
			return
		}
		
		if text == t {
			tmp++
		}
		//For debugging. To show only Monday.
//		if tmp > 2 {
//			return
//		}
		
		if tmp == 3 {
			fmt.Println("====================================")
			fmt.Println(tmp, text)
			fmt.Println("====================================")
			putToDB(departments, db)
			departments = clean(departments)
			tmp = 1
			subjectIndex = -1
		}

		if strings.Contains(class, tdtime) {
			if time == "" {
				fmt.Println("====if =============", subjectIndex, text, "=================")
				time = text
				nextLine = false
			} else if time == text {
				fmt.Println("====else if =============", subjectIndex, text, "=================")
				nextLine = true
				spanIndex = 0
				ind = 0
//				numberBeforeSmall0 = 0
			} else {
				fmt.Println("== else ===============", subjectIndex, text, "=================")
				Spans = make([]Interval, 10, 10)
				nextLine = false
				time = text
				subjectIndex++
				spanIndex = 0
//				numberBeforeSmall0 = 0
				ind = 0
			}
		}

		//To count all small0 classes
		std.Find("td").Each(func(i int, sel *goquery.Selection) {
			if small, ok := sel.Attr("class"); ok {
				if strings.Contains(small, "tdsmall0") {
					isFirstInSmall0 = true
					countSmall0++
				}
			}
		})

		if countSmall0 <= 0 {
			insertedGroups = make([]string, 5)
		}

		if countSmall0 > 0 && class != tdsmall+"0" {
			numberBeforeSmall0 = fromStringToInt(class)
			classBeforeSmall0 = class
			return
		} else if countSmall0 == 0 {
			numberBeforeSmall0 = 0
			classBeforeSmall0 = class
		}
		
		var allGr = make([]string, 0, 5)
		var room string
		std.Find("nobr").Each(func(i int, sel *goquery.Selection) {
			room = sel.Text()
		})

		if strings.Contains(classBeforeSmall0, tditem) && countSmall0 > 0 {
			is2Weeks = true
		} else {
			is2Weeks = false
		}

		if strings.Contains(class, tditem) {
			numberFromClass := fromStringToInt(class)
			subject := parseGroups(text, room)
			for i := ind; i < ind+numberFromClass; i++ {
				allGr = append(allGr, eachColumn[i]...)
			}
			
			var withoutGroups = make([]string, 1)
			var mapWithoutGr = make(map[string]string)
			var resFromReg = make([]string, 0, 5)
			if strings.Contains(text, "без"){
				withoutGroups = reGrp.FindAllString(text, -1)
				for _, v := range withoutGroups{
					mapWithoutGr[v] = ""
				}
				for _, v1 := range allGr {
					if _, ok := mapWithoutGr[v1]; ok{
						continue
					}
					resFromReg = append(resFromReg, v1)
				}
			} else {
				resFromReg = reGrp.FindAllString(text, -1)
			}
			
			st := DataToParsingLine{
				Departments:      departments,
				AllGroups:        allGr,
				ResultFromReqexp: resFromReg,
				InsertedGroups:   insertedGroups,
				Lesson:           subject,
				RegexpInterval:   reInterval,
			}
			departments, insertedGroups = st.parseLine(subjectIndex, countSmall0-1, text, nextLine, is2Weeks, isFirstInSmall0)
			ind = ind + numberFromClass

		} else if strings.Contains(class, tdsmall) {
			numberFromClass := fromStringToInt(class)
			subject := parseGroups(text, room)
			
//			resFromReg := reGrp.FindAllString(text, -1)

			if numberBeforeSmall0 == 0 {
				numberBeforeSmall0 = numberFromClass
			}

			if !nextLine {
				if !strings.Contains(class, tdsmall+"0") || !strings.Contains(classBeforeSmall0, tditem) {
					if spanIndex == 0 || (Spans[spanIndex-1].Start != ind && Spans[spanIndex-1].End != ind+numberBeforeSmall0) {
						span := Interval{Start: ind, End: ind + numberBeforeSmall0}
						Spans[spanIndex] = span
						fmt.Println("SPANS!!!!!", spanIndex, Spans[spanIndex])
						fmt.Println("numberBeforeSmall0", numberBeforeSmall0)
						spanIndex++
					}
				}
				for i := ind; i < ind+numberBeforeSmall0; i++ {
					allGr = append(allGr, eachColumn[i]...)
				}
			} else { //NEXT STRING
//				for _, v := range departments {
//					fmt.Println(v)
//				}
				is2Weeks = false
				var End int
				if numberBeforeSmall0 < Spans[spanIndex].End - Spans[spanIndex].Start {
					length := len(Spans) - 1
					isChanged := false
					Spans = append(Spans, Interval{Start: 0, End: 0})
					for j := length; j >= 0; j--{
						if j == spanIndex {
							Spans[j+1] = Interval{Start: Spans[j].Start + 1, End: Spans[j].End}
							Spans[j] = Interval{Start: Spans[j].Start, End: Spans[j].Start + 1}
							isChanged = true
							continue
						}
						if isChanged {
							continue
						}
						Spans[j+1] = Spans[j]
					}
					End = Spans[spanIndex].End
				} else if numberBeforeSmall0 == 1{
						End = Spans[spanIndex].End
				} else {
					End = numberBeforeSmall0
				}
				
				fmt.Println("!!!!!!!!NEXT STRING", Spans[spanIndex], spanIndex,"NBS", numberBeforeSmall0)
				for i := Spans[spanIndex].Start; i < End; i++ {
					allGr = append(allGr, eachColumn[i]...)
				}
				if countSmall0-1 <= 0 {
					spanIndex++
				}
			}
			
			var withoutGroups = make([]string, 1)
			var mapWithoutGr = make(map[string]string)
			var resFromReg = make([]string, 0, 5)
			if strings.Contains(text, "без"){
				withoutGroups = reGrp.FindAllString(text, -1)
				for _, v := range withoutGroups{
					mapWithoutGr[v] = ""
				}
				for _, v1 := range allGr {
					if _, ok := mapWithoutGr[v1]; ok{
						continue
					}
					resFromReg = append(resFromReg, v1)
				}
			} else {
				resFromReg = reGrp.FindAllString(text, -1)
			}
			
			st := DataToParsingLine{
				Departments:      departments,
				AllGroups:        allGr,
				ResultFromReqexp: resFromReg,
				InsertedGroups:   insertedGroups,
				Lesson:           subject,
				RegexpInterval:   reInterval,
			}
			fmt.Println("++++++++++=",text,  resFromReg)
			departments, insertedGroups = st.parseLine(subjectIndex, countSmall0-1, text, nextLine, is2Weeks, isFirstInSmall0)
			isFirstInSmall0 = false
			//very strange part...
			if countSmall0 > 0 {
				countSmall0--
				if countSmall0 != 0 {
					return
				}
				ind = ind + numberBeforeSmall0
				numberBeforeSmall0 = 0
				return
			}
			ind = ind + numberFromClass
		}
	})
}

