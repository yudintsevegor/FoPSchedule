package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"log"
)

func main(){
	res, err := http.Get("http://ras.phys.msu.ru/table/4/1.html")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	
	if res.StatusCode != 200 {
		log.Fatal("status code error: %d %s", res.StatusCode, res.Status)
	}
	
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil{
		log.Fatal(err)
	}
	
	doc.Find("a").Each(func(i int, s *goquery.Selection){
		link, _ := s.Attr("href")
		fmt.Println(link)
	})
}
