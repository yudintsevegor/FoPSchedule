package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"regexp"
	"strconv"
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

	var flag bool
	groups := make(map[string]string)
	departments := make([]Department, 0, 5)
	test := make(map[string]string)

	eachColumn := make(map[int][]string)

	indx := 0
	doc.Find("td").Each(func(i int, std *goquery.Selection) {
		//		fmt.Println("TD")
		if class, ok := std.Attr("class"); ok {
			test[class] = ""
		}

		if grpsFound > 1 {
			return
		}
		depart := Department{}
		text := std.Text()
		if flag && text != grpEnd {
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
			flag = true
		} else if text == grpEnd {
			flag = false
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
	
	doc.Find("td").Each(func(i int, std *goquery.Selection) {
		//		fmt.Println("TD")
		text := std.Text()

		if class, ok := std.Attr("class"); ok {
			if text == t{
				tmp++
			}
			if tmp > 2{
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

			if strings.Contains(class, tditem) {
				number, err := fromStringToInt(class)
				if err != nil {
					log.Fatal(err)
				}
//				depart := parseGroups(text)
				
				//				fmt.Println(class, number, len(binaryColumns), quantity)
				for i := ind; i < ind+number; i++ {
					binaryColumns[i]++
				}
				ind = ind + number
			} else if strings.Contains(class, tdsmall) {
				number, err := fromStringToInt(class)
				if err != nil {
					log.Fatal(err)
				}

				if number == 0 {
					return
				}
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
					//					fmt.Println(class, ind, number, len(binaryColumns), quantity)
					//						k := 0
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

func parseGroups(text string, depart Department) Department {
	subj := Subject{}
	
	if text == ""{
//		subj.Name = ""
//		subj.Lector = ""
//		subj.Room = ""
		depart.Lessons = append(depart.Lessons, subj)
		
		return depart
	}
	if strings.Contains(text, practice){
		subj.Name = practice
		depart.Lessons = append(depart.Lessons, subj)
		
		return depart
	}

	rSubj := regexp.MustCompile(depart.Number + ` - ([^0-9||Каф.]*)`)
	resSbj := rSubj.FindStringSubmatch(text)
	Subj := resSbj[1]
	
	rRoom := regexp.MustCompile(resSbj[0] + ` ([\d+\-\d+||Каф.]*)`)
	resRm := rRoom.FindStringSubmatch(text)
	Room := resRm[1]

	rLect := regexp.MustCompile(resRm[0] + ` (.+)`)
	Lect := rLect.FindStringSubmatch(text)[1]
	
	subj.Name = Subj
	subj.Lector = Lect
	subj.Room = Room
	
	depart.Lessons = append(depart.Lessons, subj)

	return depart
}
/**/









