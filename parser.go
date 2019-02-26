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
	grpend := "<< ГРУППЫ"

	tditem := "tditem"
	tdsmall := "tdsmall"
	tdtime := "tdtime"

	var flag bool
	groups := make(map[string]string)
	test := make(map[string]string)

	doc.Find("td").Each(func(i int, std *goquery.Selection) {
		//		fmt.Println("TD")
		if class, ok := std.Attr("class"); ok {
			test[class] = ""

		}

		text := std.Text()
		if flag && text != grpend {
			groups = parseGroups(text, groups)
			groups[text] = ""
		}
		if text == grpbegin {
			flag = true
		} else if text == grpend {
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
	var time string
	var nextStr bool

	re := regexp.MustCompile(`[a-zA-z]([0-9]+)`)
	var ind int
	
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
				} else {
					time = text
					nextStr = true
					ind = 0
//					fmt.Println(time, binaryColumns)
					for i := 0; i < quantity - 1; i++{
						binaryColumns[i] = 0
					}
				}
			}

			if strings.Contains(class, tditem) {
				num := re.FindStringSubmatch(class)[1]
				number, err := strconv.Atoi(num)
				if err != nil {
					log.Fatal(err)
				}
//				fmt.Println(class, number, len(binaryColumns), quantity)
				for i := ind; i < ind + number - 1; i++{
					binaryColumns[i]++
				}
				ind = ind + number
			} else if strings.Contains(class, tdsmall) {
				num := re.FindStringSubmatch(class)[1]
				number, err := strconv.Atoi(num)
				if err != nil {
					log.Fatal(err)
				}
				if number != 0 {
	//			fmt.Println(class, ind, number, len(binaryColumns), quantity)
				for i := ind; i < ind + number; i++{
					binaryColumns[i] += 0.5
				}
				
				fmt.Println(ind, ind + number - 1, time, binaryColumns)
				ind = ind + number
				}
			}
		}
		//		fmt.Println(text)
	})
}

func parseGroups(text string, mapka map[string]string) map[string]string {
	return mapka
}
