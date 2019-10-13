package main

import (
	"fmt"
	"net/http"
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
	case cookieURL:
		h.handleCookie(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	sessions := make(map[string]User)
	handle := &Handler{
		Sessions: sessions,
	}

	fmt.Println("starting server at :" + port)
	http.ListenAndServe(":"+port, handle)
}
