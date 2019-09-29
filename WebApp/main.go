package main

import (
	"fmt"
	"net/http"
)

const (
	host = "http://localhost:8080"
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
	case "/cookie":
		h.handleCookie(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	var sessions = make(map[string]User)
	handle := &Handler{
		Sessions: sessions,
	}
	port := "8080"
	fmt.Println("starting server at :" + port)
	http.ListenAndServe(":"+port, handle)
}
