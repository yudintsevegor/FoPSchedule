package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	Tmpl *template.Template
}

type PriceStat struct {
	Time   int
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int
	Ticker string
}

type Transmit struct {
	ID int
	UserID int `json:"user_id"`
	Vol    int
	Price  int
	IsBuy  int `json:"is_buy"`
	Ticker string
}

type Test struct{
	Group string `json:"group"`
	Code string `json:"code"`
}

func (h *Handler) Transcation(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    fmt.Fprintln(w, r.Form)
	json := r.Form[data][0]
	test := Test{}
	err := json.Unmarshal(json, &test)
	if err != nil{
		log.Fatal(err)
	}
	
}

func (h *Handler) Position(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8082/status", nil)
	if err != nil {
		fmt.Println(err)
	}
	client := &http.Client{Timeout: time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(body)
	var statistic map[string][]*Transmit
	err = json.Unmarshal(body, &statistic)
	if err != nil {
		fmt.Println(err)
	}
	res := statistic["result"]
	err = h.Tmpl.ExecuteTemplate(w, "position.html", struct {
		Stats []*Transmit
	}{
		Stats: res,
	})
	if err != nil {
		fmt.Println(err)
	}
}

func (h *Handler) Price(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8082/stat", nil)
	if err != nil {
		fmt.Println(err)
	}
	client := &http.Client{Timeout: 5*time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var statistic map[string][]*PriceStat
	err = json.Unmarshal(body, &statistic)
	if err != nil {
		fmt.Println(err)
	}
	res := statistic["result"]
	err = h.Tmpl.ExecuteTemplate(w, "price.html", struct {
		Stats []*PriceStat
	}{
		Stats: res,
	})
	if err != nil {
		fmt.Println(err)
	}
}

func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Продукт находится в активной разработке. По вопросам обращаться по номеру +4 8 15 16 23 42"))
	info := &Transmit{
		ID: 1, //NEED TO FIX, WANT TO GET ID!!!!!!!!!
	}
	res, err := json.Marshal(info)
	if err != nil {
		fmt.Println(err)
	}
	reqBody := bytes.NewReader(res)
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8082/cancel", reqBody)
	client := &http.Client{Timeout: time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var statistic map[string][]*Transmit
	err = json.Unmarshal(body, &statistic)
	if err != nil {
		fmt.Println(err)
	}
	result := statistic["result"]
	err = h.Tmpl.ExecuteTemplate(w, "position.html", struct {
		Stats []*Transmit
	}{
		Stats: result,
	})

	if err != nil {
		fmt.Println(err)
	}

}

func (h *Handler) Action(w http.ResponseWriter, r *http.Request) {
	sell := r.FormValue("sell")
	buy := 1
	if sell == "Продать" {
		buy = 0
	}
	volume, err := strconv.Atoi(r.FormValue("vol"))
	if err != nil {
		fmt.Println(err)
	}
	price, err := strconv.Atoi(r.FormValue("price"))
	if err != nil {
		fmt.Println(err)
	}
	info := &Transmit{
		UserID: 100500,
		Vol:    volume,
		Price:  price,
		IsBuy:  buy,
		Ticker: r.FormValue("ticker"),
	}
	res, err := json.Marshal(info)
	if err != nil {
		fmt.Println(err)
	}
	reqBody := bytes.NewReader(res)
	req, err := http.NewRequest(http.MethodPut, "http://127.0.0.1:8082/deal", reqBody)

	client := &http.Client{Timeout: time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r,"/", http.StatusFound)
	return
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "test.html", struct{}{})
	if err != nil {
		fmt.Println(err)
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		h.List(w, r)
	case "/position":
		h.Position(w, r)
	case "/price":
		h.Price(w, r)
	case "/form":
		h.Action(w, r)
	case "/cancel":
		h.Cancel(w, r)
	case "/result":
		h.Transcation(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	handler := &Handler{
		Tmpl: template.Must(template.ParseGlob("templates/*")),
	}
	fmt.Println("starting server at :8080")
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		fmt.Println(err)
	}
}
