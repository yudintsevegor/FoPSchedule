package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var config *oauth2.Config

func init() {
	config = &oauth2.Config{
		RedirectURL:  host + cookieURL,
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/calendar", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func (h *Handler) handleMain(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob(mainHTMLPage)
	if err != nil {
		log.Fatal(err)
	}
	tmpl.ExecuteTemplate(w, mainHTMLPage, nil)
}

func (h *Handler) handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	h.Mutex.Lock()
	oauthStateString := getRandomString()
	url := config.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline)
	h.Mutex.Unlock()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) handleCookie(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie(cookieName)
	if err == nil {
		for _, c := range r.Cookies() {
			if c.Name == cookieName {
				http.SetCookie(w, &http.Cookie{
					Name:    c.Name,
					MaxAge:  -1,
					Expires: time.Now().Add(-100 * time.Minute),
				})

				if _, ok := h.Sessions[c.Value]; ok {
					h.Mutex.Lock()
					delete(h.Sessions, c.Value)
					h.Mutex.Unlock()
				}
			}
		}
	}

	oauthStateString := getRandomString()
	cook := &http.Cookie{
		Name:  cookieName,
		Value: oauthStateString,
	}

	http.SetCookie(w, cook)
	r.AddCookie(cook)
	c, err := r.Cookie(cookieName)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("added cookie: ", c)
	}

	code := r.FormValue("code")
	client, email, err := getClient(code)
	if err != nil {
		log.Fatal(err)
	}

	h.Mutex.Lock()
	st := h.Sessions[cook.Value]
	st.Client = client
	st.Email = email
	h.Sessions[cook.Value] = st
	h.Mutex.Unlock()

	http.Redirect(w, r, host+"/callback", http.StatusSeeOther)
}

func (h *Handler) handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(cookieName)
	if err != nil {
		fmt.Println("Error in '/callback': ", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("Error, " + err.Error()))
		http.Redirect(w, r, host+"/login", http.StatusSeeOther)
		return
	}

	if _, ok := h.Sessions[c.Value]; !ok {
		fmt.Println("no value at map")
		http.Redirect(w, r, host+"/login", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseGlob("index.html")
	if err != nil {
		log.Fatal(err)
	}
	h.Mutex.Lock()
	email := h.Sessions[c.Value].Email
	h.Mutex.Unlock()
	path := host + "/result"
	tmpl.ExecuteTemplate(w, "index.html", User{Email: email, PathAction: path})
}

func (h *Handler) handleResult(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(cookieName)
	if err != nil {
		fmt.Fprintf(w, "no cookie")
		return
	}

	h.Mutex.Lock()
	_, ok := h.Sessions[c.Value]
	h.Mutex.Unlock()
	if !ok {
		http.Redirect(w, r, host+"/login", http.StatusTemporaryRedirect)
		return
	}

	group := r.FormValue("group")
	h.Mutex.Lock()
	client := h.Sessions[c.Value].Client
	h.Mutex.Unlock()

	go putData(client, group)
	http.Redirect(w, r, urlCalendar, http.StatusTemporaryRedirect)
}

func getClient(code string) (*http.Client, string, error) {
	client := &http.Client{}
	token, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return client, "", fmt.Errorf("code exchange failed: %s", err.Error())
	}
	client = config.Client(oauth2.NoContext, token)
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	info := UserInfo{}
	err = json.Unmarshal(contents, &info)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(info)

	return client, info.Email, nil
}

func getRandomString() string {
	size := 16
	rb := make([]byte, size)
	_, err := rand.Read(rb)
	if err != nil {
		log.Fatal(err)
	}
	oauthStateString := base64.URLEncoding.EncodeToString(rb)

	return oauthStateString
}
