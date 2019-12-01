package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"

	"fopSchedule/master/common"
)

const (
	port = "8080"
	host = "http://localhost:" + port
)

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		h.handleMain(w, r)
	case "/login":
		h.handleGoogleLogin(w, r)
	case "/callback":
		h.handleGoogleCallback(w, r)
	case "/result":
		h.handleResult(w, r)
	case common.CookieURL:
		h.handleCookie(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	db, err := sql.Open("mysql", DSN)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxIdleConns(common.MaxConnections)

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	sessions := make(map[string]User)
	handle := &Handler{
		Sessions: sessions,
		Mutex:    &sync.Mutex{},
		DB:       db,
	}

	fmt.Println("starting server at :" + port)
	http.ListenAndServe(":"+port, handle)
}
