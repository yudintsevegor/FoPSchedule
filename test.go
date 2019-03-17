package main

import (
	"fmt"
	"net/http"
//	"io/ioutil"
	"log"
	"os"
)



func main() {

	var url string
	fmt.Fscan(os.Stdin, &url)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println(req.FormValue("approvalCode"))
//	resp, err := http.Get(url)
//	
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(string(body))
}
