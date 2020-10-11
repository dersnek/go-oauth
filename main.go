package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

const (
	loginURL         = "https://api.github.com/user"
	oauthStateString = "pseudo-random"
)

var (
	conf = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     "CLIENT_ID",
		ClientSecret: "CLIENT_SECRET",
		Endpoint:     endpoint,
	}
	endpoint = oauth2.Endpoint{
		AuthURL:  "https://github.com/login/oauth/authorize",
		TokenURL: "https://github.com/login/oauth/access_token",
	}
)

// User - simple user
type User struct {
	Login string `json:"login"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleCallback)
	http.ListenAndServe(":8080", nil)
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
		<html><body><a href="/login">Sign in with GitHub</a></body></html>
	`)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	url := conf.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	token, err := conf.Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, loginURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "token "+token.AccessToken)

	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var u User
	err = json.Unmarshal(content, &u)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("User: %+v\n", u)
	w.Write(content)
}
