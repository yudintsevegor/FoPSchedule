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

	grpbegin := "ГРУППЫ >>"
	grpEnd := "<< ГРУППЫ"

	var flag bool
	groups := make(map[string]string)
	test := make(map[string]string)

	doc.Find("td").Each(func(i int, std *goquery.Selection) {
		//		fmt.Println("TD")
		if class, ok := std.Attr("class"); ok {
			test[class] = ""

		}

		text := std.Text()
		if flag && text != grpEnd {
			groups = parseGroups(text, groups)
			groups[text] = ""
		}
		if text == grpbegin {
			flag = true
		} else if text == grpEnd {
			flag = false
		}
		//		fmt.Println(text)
	})

	for key, _ := range groups {
		fmt.Println(key)
	}

	quantity := len(groups)
	//	fmt.Println(quantity)
	var binaryColumns = make([]float64, quantity)
//	var indexSlice = make([]int, quantity)
	
	var time string
	var nextStr bool

	var ind int
	tditem := "tditem"
	tdsmall := "tdsmall"
	tdtime := "tdtime"

	doc.Find("td").Each(func(i int, std *goquery.Selection) {
		//		fmt.Println("TD")
		text := std.Text()

		if class, ok := std.Attr("class"); ok {
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

				if number != 0 {
//					fmt.Println(nextStr)
					if !nextStr {
//						for j := 0; j < number; j++{
//							binaryColumns[indexSlice[j]] += 0.5
//						}
//						indexSlice = indexSlice[:number]
						for j := 0; j < number; j++{
							for i := 0; i < quantity; i++ {
								if binaryColumns[i] == 0.5 {
									binaryColumns[i] += 0.5
									break
								}
							}
						}
					} else {
//					fmt.Println(class, ind, number, len(binaryColumns), quantity)
						k := 0
						for i := ind; i < ind+number; i++ {
//							fmt.Println(k)
//							indexSlice[k] = i
							if binaryColumns[i] == 0 {
								binaryColumns[i] += 0.5
							}
							k++
						}
						ind = ind + number
					}
				}
			}
			fmt.Println(time, binaryColumns, class)
		}
		//		fmt.Println(text)
	})
}

func fromStringToInt(class string) (int, error) {
	num := re.FindStringSubmatch(class)[1]
	number, err := strconv.Atoi(num)

	return number, err
}

func parseGroups(text string, mapka map[string]string) map[string]string {
	return mapka
}
