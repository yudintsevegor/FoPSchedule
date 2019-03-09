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
	res, err := http.Get("http://ras.phys.msu.ru/table/4/2.html")
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
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if link, ok := s.Attr("href"); ok {
			fmt.Println(link)
			text := s.Text()
			fmt.Println(text)
		}
	})

	course := "4"
	var reGrp = regexp.MustCompile(course + `\d{2}`)

	grpBegin := "ГРУППЫ >>"
	grpEnd := "<< ГРУППЫ"
	var grpsFound int

	var isGroups bool
	groups := make(map[string]string)
	var departments = make([]Department, 0, 5)
	Clss := make(map[string]string)

	eachColumn := make(map[int][]string)

	indx := 0
	doc.Find("td").Each(func(i int, std *goquery.Selection) {
		if class, ok := std.Attr("class"); ok {
			Clss[class] = ""
		}

		if grpsFound > 1 {
			return
		}
		//if set []Subject, 0, 5, program will panic. WHY?

		text := std.Text()
		if isGroups && text != grpEnd {
			resFromReg := reGrp.FindAllString(text, -1)
			eachColumn[indx] = resFromReg
			indx++
			for _, val := range resFromReg {
				depart := Department{Lessons: make([]Subject, 5, 5)}
				d := Department{Lessons: make([]Subject, 5, 5)}
				depart.Number = val
				d.Number = val
				departments = append(departments, depart)
			}
			groups[text] = ""
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

	db, err := sql.Open("mysql", DSN)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	partOfReq := `(
				  id int(11) NOT NULL AUTO_INCREMENT,
				  first varchar(255),
				  second varchar(255),
				  third varchar(255),
				  fourth varchar(255),
				  fifth varchar(255),
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
	var nextStr bool
	var ind int
	tditem := "tditem"
	tdsmall := "tdsmall"
	tdtime := "tdtime"
	t := "9:00- - -  10:35"
	var tmp int
	var classBeforeSmall0 string
	var numberBeforeSmall0 int
	var countSmall0 int
	var n int
	var is2Weeks bool
	var Spans = make([]Interval, 10, 10)
	var insertedGroups = make([]string, 5)
	var num int
	var Saturday int
	var columns = " ( first, second, third, fourth, fifth ) "
	var quesStr = " ( ?, ?, ?, ?, ? ) "

	isSat := false
	doc.Find("td").Each(func(i int, std *goquery.Selection) {
		text := std.Text()

		if text == grpBegin {
			Saturday++
			fmt.Println("================================================", text, Saturday, "==================================")
		}

		if Saturday == 3 && !isSat {
			for _, val := range departments {
				var valuesToDB = make([]interface{}, 0, 1)
				fmt.Println(val.Number)
				for _, les := range val.Lessons {
					value := les.Name + "%" + les.Lector + "%" + les.Room
					valuesToDB = append(valuesToDB, value)
				}
				fmt.Println(valuesToDB)
				req := fmt.Sprintf("INSERT INTO `%v`"+columns+"VALUES"+quesStr, val.Number)
				statement, err := db.Prepare(req)
				if err != nil {
					log.Fatal(err)
				}
				_, err = statement.Exec(valuesToDB...)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("PUT IN TABLE")
			}
			isSat = true
		}

		if class, ok := std.Attr("class"); ok {
			if text == t {
				tmp++
			}
			//For debugging. To show only Monday.
			//			if tmp > 2 {
			//				return
			//			}

			//			if tmp > 6 || tmp < 5 {
			if tmp == 3 {
				fmt.Println("====================================")
				fmt.Println(tmp, text)
				fmt.Println("====================================")
				for _, val := range departments {
					var valuesToDB = make([]interface{}, 0, 1)
					fmt.Println(val.Number)
					for _, les := range val.Lessons {
						value := les.Name + "%" + les.Lector + "%" + les.Room
						valuesToDB = append(valuesToDB, value)
					}
					fmt.Println(valuesToDB)
					req := fmt.Sprintf("INSERT INTO `%v`"+columns+"VALUES"+quesStr, val.Number)
					statement, err := db.Prepare(req)
					if err != nil {
						log.Fatal(err)
					}
					_, err = statement.Exec(valuesToDB...)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Println("PUT IN TABLE")
				}
				//				fmt.Println(emptyDep, tmp)
				tmp = 1
				//				departments = emptyDep
				n = -1
				fmt.Println(n)
				//				fmt.Println(departments, tmp)
				departments = clean(departments)
				//				fmt.Println(departments, tmp)
			}

			if strings.Contains(class, tdtime) {
				if time == "" {
					fmt.Println("====if =============", n, text, "=================")
					time = text
					nextStr = false
				} else if time == text {
					fmt.Println("====else if =============", n, text, "=================")
					num = 0
					nextStr = true
					ind = 0
					numberBeforeSmall0 = 0
				} else {
					fmt.Println("== else ===============", n, text, "=================")
					num = 0
					Spans = make([]Interval, 10, 10)
					n++
					time = text
					nextStr = false
					ind = 0
				}
			}

			std.Find("td").Each(func(i int, sel *goquery.Selection) {
				if small, ok := sel.Attr("class"); ok {
					if strings.Contains(small, "tdsmall0") {
						countSmall0++
					}
				}
			})

			if countSmall0 <= 0 {
				insertedGroups = make([]string, 5)
			}

			if countSmall0 > 0 && class != tdsmall+"0" {
				number := fromStringToInt(class)
				numberBeforeSmall0 = number
				classBeforeSmall0 = class
				return
			} else if countSmall0 == 0 {
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

			//			if !(strings.Contains(class, tditem) && strings.Contains(class, tdsmall){
			//				return
			//			}
			if strings.Contains(class, tditem) {
				number := fromStringToInt(class)
				subject := parseGroups(text, room)
				resFromReg := reGrp.FindAllString(text, -1)

				for i := ind; i < ind+number; i++ {
					allGr = append(allGr, eachColumn[i]...)
				}

				departments, insertedGroups = parseLine(departments, allGr, resFromReg, insertedGroups, subject, text, countSmall0-1, n, nextStr, is2Weeks)
				ind = ind + number

			} else if strings.Contains(class, tdsmall) {
				number := fromStringToInt(class)
				subject := parseGroups(text, room)
				resFromReg := reGrp.FindAllString(text, -1)

				if numberBeforeSmall0 == 0 {
					numberBeforeSmall0 = number
				}

				if !nextStr {
					//					fmt.Println("==========================inserted groups and countSmall0 and is2Weeks================================")
					//					fmt.Println(insertedGroups, countSmall0, is2Weeks)
					//					fmt.Println("==========================================================")

					//					fmt.Println(class, ind, numberBeforeSmall0, classBeforeSmall0)
					if !strings.Contains(class, tdsmall+"0") || !strings.Contains(classBeforeSmall0, tditem) {
						//						fmt.Println("SPANS before", num, Spans[num])
						if num == 0 || (Spans[num-1].Start != ind && Spans[num-1].End != ind+numberBeforeSmall0) {
							span := Interval{Start: ind, End: ind + numberBeforeSmall0}
							Spans[num] = span
							fmt.Println("SPANS!!!!!", num, Spans[num])
							num++
						}
					}
					for i := ind; i < ind+numberBeforeSmall0; i++ {
						allGr = append(allGr, eachColumn[i]...)
					}
					//departments = parseLine(departments, allGr, resFromReg, insertedGroups, subject, text, n, nextStr)
				} else { //NEXT STRING
					fmt.Println("LOL", countSmall0)
					for _, v := range departments {
						fmt.Println(v)
					}
					is2Weeks = false
					fmt.Println(Spans[num], num, text)
					for i := Spans[num].Start; i < Spans[num].End; i++ {
						allGr = append(allGr, eachColumn[i]...)
					}
					//departments = parseLine(departments, allGr, resFromReg, insertedGroups, subject, text, countSmall0, n, nextStr)
					if countSmall0-1 <= 0 {
						num++
					}
				}
				departments, insertedGroups = parseLine(departments, allGr, resFromReg, insertedGroups, subject, text, countSmall0-1, n, nextStr, is2Weeks)

				if countSmall0 > 0 {
					countSmall0--
					if countSmall0 != 0 {
						return
					}
					ind = ind + numberBeforeSmall0
					numberBeforeSmall0 = 0
					return
				}
				ind = ind + number
			}

			//			fmt.Println(ind, time, class, text, "\n")
		}
	})

}
