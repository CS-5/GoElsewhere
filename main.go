package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type (
	Config struct {
		DefaultURL string     `json:"default"`
		Port       string     `json:"port"`
		Redirects  []Redirect `json:"redirects"`
	}

	Redirect struct {
		Code string `json:"code"`
		URL  string `json:"url"`
	}
)

var config Config

func main() {
	file, fileErr := ioutil.ReadFile("config.json")
	if fileErr != nil {
		fmt.Println(fileErr)
	}

	err := json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Starting HTTP Server")
	http.HandleFunc("/", redirectHandler)

	http.ListenAndServe(":"+config.Port, nil)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request recieved for: " + r.URL.String())

	//Check for default requests
	if r.URL.String() == "/" || r.URL.String() == "/favicon.ico" {
		fmt.Println("Sending user to default")
		http.Redirect(w, r, config.DefaultURL, http.StatusSeeOther)
	} else {

		//Iterate through redirects
		for i := range config.Redirects {

			//If redirect is found, send the user
			if config.Redirects[i].Code == r.URL.String() {
				fmt.Println("Sending user to: " + config.Redirects[i].URL)
				http.Redirect(w, r, config.Redirects[i].URL, http.StatusSeeOther)
			}
		}

		//If for loop exits without finding anything, send user to default    <--- Is there a better way to do this?
		fmt.Println("Sending user to default")
		http.Redirect(w, r, config.DefaultURL, http.StatusSeeOther)
	}
}
