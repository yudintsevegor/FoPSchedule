package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"fopSchedule/master/common"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var config *oauth2.Config

func init() {
	config = &oauth2.Config{
		RedirectURL:  host + common.CookieURL,
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/calendar", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func (h *Handler) handleMain(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob(common.MainHTMLPage)
	if err != nil {
		log.Fatal(err)
		return
	}
	tmpl.ExecuteTemplate(w, common.MainHTMLPage, nil)
}

func (h *Handler) handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	oauthStateString := getRandomString()
	url := config.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) handleCookie(w http.ResponseWriter, r *http.Request) {
	cookieName := common.CookieName
	_, err := r.Cookie(cookieName)
	if err == nil {
		for _, c := range r.Cookies() {
			if c.Name != cookieName {
				continue
			}

			http.SetCookie(w, &http.Cookie{
				Name:    c.Name,
				MaxAge:  -1,
				Expires: time.Now().Add(-100 * time.Minute),
			})

			h.Mutex.Lock()
			if _, ok := h.Sessions[c.Value]; ok {
				delete(h.Sessions, c.Value)
			}
			h.Mutex.Unlock()
		}
	}

	oauthStateString := getRandomString()
	cook := &http.Cookie{
		Name:  cookieName,
		Value: oauthStateString,
	}

	http.SetCookie(w, cook)
	r.AddCookie(cook)
	_, err = r.Cookie(cookieName)
	if err != nil {
		log.Println(err)
		return
	}

	client, email, err := getClientAndInfo(r.FormValue("code"))
	if err != nil {
		log.Println(err)
		return
	}

	h.Mutex.Lock()
	h.Sessions[cook.Value] = User{
		Client: client,
		Email:  email,
	}
	h.Mutex.Unlock()

	http.Redirect(w, r, host+"/callback", http.StatusSeeOther)
}

func getClientAndInfo(code string) (*http.Client, string, error) {
	if code == "" {
		return &http.Client{}, "", errors.New("code is empty")
	}

	token, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return &http.Client{}, "", fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return &http.Client{}, "", err
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &http.Client{}, "", err
	}

	info := UserInfo{}
	if err = json.Unmarshal(contents, &info); err != nil {
		return &http.Client{}, "", err
	}

	client := config.Client(oauth2.NoContext, token)

	return client, info.Email, nil
}

func (h *Handler) handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(common.CookieName)
	if err != nil {
		log.Println("Error in '/callback': ", err)
		errorHandler(w, http.StatusUnauthorized, err)
		// http.Redirect(w, r, host+"/login", http.StatusSeeOther)
		return
	}

	if _, ok := h.Sessions[c.Value]; !ok {
		log.Println("Error in '/callback': ", err)
		errorHandler(w, http.StatusUnauthorized, err)
		// http.Redirect(w, r, host+"/login", http.StatusSeeOther)
		return
	}

	indexHTMLPage := common.IndexHTMLPage
	tmpl, err := template.ParseGlob(indexHTMLPage)
	if err != nil {
		log.Println(err)
		return
	}

	h.Mutex.Lock()
	email := h.Sessions[c.Value].Email
	h.Mutex.Unlock()

	tmpl.ExecuteTemplate(w, indexHTMLPage, User{Email: email, PathAction: host + "/result"})
}

func (h *Handler) handleResult(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(common.CookieName)
	if err != nil {
		log.Println(err)
		errorHandler(w, http.StatusUnauthorized, err)
		return
	}

	h.Mutex.Lock()
	_, ok := h.Sessions[c.Value]
	h.Mutex.Unlock()
	if !ok {
		http.Redirect(w, r, host+"/login", http.StatusTemporaryRedirect)
		return
	}

	// bottle neck, TODO: change
	h.Mutex.Lock()
	user := h.Sessions[c.Value]
	client := user.Client
	h.Mutex.Unlock()

	group := r.FormValue("group")
	if group == "" {
		err := errors.New("doesnt match group")
		log.Println(err)
		errorHandler(w, http.StatusInternalServerError, err)
		return
	}

	go func() {
		if err := user.putData(h.DB, client, group); err != nil {
			log.Println(err)
			errorHandler(w, http.StatusInternalServerError, err)
			return
		}
	}()

	http.Redirect(w, r, host, http.StatusTemporaryRedirect)
}

func errorHandler(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Error, " + err.Error()))
	w.WriteHeader(statusCode)
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
