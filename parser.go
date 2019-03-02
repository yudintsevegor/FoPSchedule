package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode"
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

	var reGrp = regexp.MustCompile(`4\d+`)

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
		depart := Department{}
		text := std.Text()
		if isGroups && text != grpEnd {
			resFromReg := reGrp.FindAllString(text, -1)
			eachColumn[indx] = resFromReg
			indx++
			for _, val := range resFromReg {
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
		//		fmt.Println(text)
	})

	for _, key := range departments {
		fmt.Println(key)
	}

	for key, val := range eachColumn {
		fmt.Println(key, val)
	}

	quantity := len(groups)
	//	fmt.Println(quantity)
	var binaryColumns = make([]float64, quantity)

	var time string
	var nextStr bool

	var ind int
	tditem := "tditem"
	tdsmall := "tdsmall"
	tdtime := "tdtime"
	t := "9:00- - -  10:35"
	var tmp int

	var countSmall0 int
	doc.Find("td").Each(func(i int, std *goquery.Selection) {
		//		fmt.Println("TD")
		text := std.Text()

		if class, ok := std.Attr("class"); ok {
			//			For debugging. To show only Monday.
			if text == t {
				tmp++
			}
			if tmp > 2 {
				return
			}

			if strings.Contains(class, tdtime) {
				if time == text {
					//					fmt.Println(time, binaryColumns)
					nextStr = false
					ind = 0
				} else if time == "" {
					time = text
					nextStr = true
				} else {
					time = text
					nextStr = true
					ind = 0
					//					fmt.Println(time, binaryColumns)
					for i := 0; i < quantity; i++ {
						binaryColumns[i] = 0
					}
				}
			}

			std.Find("td").Each(func(i int, sel *goquery.Selection) {
				if small, ok := sel.Attr("class"); ok {
					if strings.Contains(small, "tdsmall0") {
						countSmall0++
					}
				}
			})

			if countSmall0 > 0 {
				countSmall0 = 0
				return
			}

			var room string
			std.Find("nobr").Each(func(i int, sel *goquery.Selection) {
				room = sel.Text()
			})

			if strings.Contains(class, tditem) {
				number, err := fromStringToInt(class)
				if err != nil {
					log.Fatal(err)
				}

				subject := parseGroups(text, room)
				fmt.Printf("Name: %v\nRoom: %v\nLector: %v\n\n", subject.Name, subject.Room, subject.Lector)

				//				fmt.Println(class, number, len(binaryColumns), quantity)
				for i := ind; i < ind+number; i++ {
					binaryColumns[i]++
				}
				ind = ind + number
				//				fmt.Println("-----------------------------------------------------------------------------------")
				//				rr, _ := std.Html()
				//				fmt.Println(rr)
			} else if strings.Contains(class, tdsmall) {
				number, err := fromStringToInt(class)
				if err != nil {
					log.Fatal(err)
				}

				/*
					if number == 0 {
						fmt.Println(class, text)
						return
					}
				*/
				subject := parseGroups(text, room)
				fmt.Printf("Name: %v\nRoom: %v\nLector: %v\n\n", subject.Name, subject.Room, subject.Lector)

				//				if number != 0 {
				//				fmt.Println(nextStr)
				if !nextStr {
					//						for j := 0; j < number; j++{
					//							binaryColumns[indexSlice[j]] += 0.5
					//						}
					//						indexSlice = indexSlice[:number]
					for j := 0; j < number; j++ {
						for i := 0; i < quantity; i++ {
							if binaryColumns[i] == 0.5 {
								binaryColumns[i] += 0.5
								break
							}
						}
					}
				} else {
					//	fmt.Println(class, ind, number, len(binaryColumns), quantity)
					//k := 0
					for i := ind; i < ind+number; i++ {
						//							fmt.Println(k)
						//							indexSlice[k] = i
						if binaryColumns[i] == 0 {
							binaryColumns[i] += 0.5
						}
						//							k++
					}
					ind = ind + number
				}
				//				}
			}
			fmt.Println(time, binaryColumns, class, text)
		}
		//		fmt.Println(text)
	})
}

func fromStringToInt(class string) (int, error) {
	num := re.FindStringSubmatch(class)[1]
	number, err := strconv.Atoi(num)

	return number, err
}

/**/
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

/**/
