package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"regexp"
//	"strconv"
	"strings"
)

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
	End int
}

var re = regexp.MustCompile(`[a-zA-z]([0-9]+)`)

func main() {
	res, err := http.Get("http://ras.phys.msu.ru/table/4/1.html")
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
//	var reInterval = regexp.MustCompile(`(` +  course + `\d{2})\s*\-\s*` + `(` +  course + `\d{2})`)

	grpbegin := "ГРУППЫ >>"
	grpEnd := "<< ГРУППЫ"
	var grpsFound int

	var isGroups bool
	groups := make(map[string]string)
	departments := make([]Department, 0, 5)
	Clss := make(map[string]string)

	eachColumn := make(map[int][]string)

	indx := 0
	doc.Find("td").Each(func(i int, std *goquery.Selection) {
		//		fmt.Println("TD")
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
				depart.Number = val
				departments = append(departments, depart)
			}
			groups[text] = ""
		}
		if text == grpbegin {
			grpsFound++
			isGroups = true
		} else if text == grpEnd {
			isGroups = false
		}
	})

	for key, val := range eachColumn {
		fmt.Println(key, val)
	}

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
//	var indMap = make(map[int]int)
	var Spans = make([]Interval, 10, 10)

	var num int
	doc.Find("td").Each(func(i int, std *goquery.Selection) {
		text := std.Text()

		if class, ok := std.Attr("class"); ok {
//			For debugging. To show only Monday.
			if text == t {
				tmp++
			}
//			if tmp > 6 || tmp < 5 {
			if tmp > 2  {
				return
			}

			if strings.Contains(class, tdtime) {
				if time == text {
					nextStr = true
					num = 0
					ind = 0
					numberBeforeSmall0 = 0
				} else if time == "" {
					time = text
					nextStr = false
				} else {
//					indMap = make(map[int]int)
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

			if countSmall0 > 0 && class != tdsmall + "0" {
				number, err := fromStringToInt(class)
				if err != nil {
					log.Fatal(err)
				}
				numberBeforeSmall0 = number
				classBeforeSmall0 = class
				return
			}

			var fullDay bool
			var room string
			std.Find("nobr").Each(func(i int, sel *goquery.Selection) {
				room = sel.Text()
			})

			if strings.Contains(class, tditem) {
				fmt.Println(class)
				
				number, err := fromStringToInt(class)
				if err != nil {
					log.Fatal(err)
				}
				subject := parseGroups(text, room)
				fmt.Printf("Name: %v\nRoom: %v\nLector: %v\n", subject.Name, subject.Room, subject.Lector)
				resFromReg := reGrp.FindAllString(text, -1)
				
				var allGr = make([]string, 0, 1)
				for i := ind; i < ind + number; i++{
					allGr = append(allGr, eachColumn[i]...)
				}
				
				departments = parseLine(departments, allGr, resFromReg, subject, text, n, nextStr)
				ind = ind + number
				
			} else if strings.Contains(class, tdsmall) {
				fmt.Println(class)
				number, err := fromStringToInt(class)
				if err != nil {
					log.Fatal(err)
				}
				if numberBeforeSmall0 == 0 {
					numberBeforeSmall0 = number
					fullDay = false
				} else if strings.Contains(class, tdsmall){
					fullDay = false
				} else {
					fullDay = true
				}

				subject := parseGroups(text, room)
				fmt.Printf("Name: %v\nRoom: %v\nLector: %v\n", subject.Name, subject.Room, subject.Lector)
				resFromReg := reGrp.FindAllString(text, -1)
				
				var allGr = make([]string, 0, 5)
				if !nextStr {
					fmt.Println(ind, numberBeforeSmall0, classBeforeSmall0, fullDay)
					if !strings.Contains(class, tdsmall + "0") || !strings.Contains(classBeforeSmall0, tditem) {
						if num == 0 || Spans[num].Start != ind && Spans[num].End != ind + numberBeforeSmall0 -1 {
							span := Interval{Start: ind, End: ind + numberBeforeSmall0 - 1}
							Spans[num] = span
							fmt.Println(Spans[num])
							num++
						}
					}
					for i := ind; i < ind + numberBeforeSmall0; i++{
						allGr = append(allGr, eachColumn[i]...)
					}
					departments = parseLine(departments, allGr, resFromReg, subject, text, n, nextStr)

				 } else { //NEXT STRING
					for i := Spans[num].Start; i < Spans[num].End + 1; i++{
						allGr = append(allGr, eachColumn[i]...)
					}
					departments = parseLine(departments, allGr, resFromReg, subject, text, n, nextStr)
					num++
				}
				
//				for _, dep := range departments{
//						fmt.Println(dep)
//				}
				
				if countSmall0 > 0{
					countSmall0--
					if countSmall0 != 0{
						return
					}
					ind = ind + numberBeforeSmall0
					numberBeforeSmall0 = 0
					return
				}
				ind = ind + number
			}
			fmt.Println(ind, time, class, text, "\n")
		}
	})
	for _, val := range departments{
		fmt.Println(val.Number)
		fmt.Println(val.Lessons, "\n")
	}
}



//						countSmall0--
//						if countSmall0 == 0{
//							fmt.Println("===================================================")
//							fmt.Println(allGr)
//							fmt.Println(resFromReg)
//							grWithEmpty := make([]string, 0, 1)
//							for i, a := range allGr{
//								for _, b := range resFromReg{
//									if a == b{
//										allGr[i] = "0"
//									}
//								}
//							}
//							for _, a := range allGr{
//								if a != "0"{
//									grWithEmpty = append(grWithEmpty, a)
//								}
//							}
//							fmt.Println(grWithEmpty)
//							for _, dep := range departments{
//								for _, gr := range grWithEmpty{
//									if dep.Number == gr {
//										newSubj := Subject{}
//										newSubj.Name = "__"
//										dep.Lessons[n] = newSubj
//									}
//								}
//							}
//						}
//						countSmall0++

