package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	RESOURCE_URL     = "http://localhost:5555"
	AUTH_SERVICE_URL = "http://localhost:8080"
)

func testSendToken() {
	fmt.Println("HTTP JSON POST URL:", RESOURCE_URL)

	var jsonData = []byte(`{
		"token": "12345"
	}`)
	request, error := http.NewRequest("POST", RESOURCE_URL, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	log.Println("Sending auth request...")
	buf, _ := ioutil.ReadAll(r.Body)
	rdr := ioutil.NopCloser(bytes.NewBuffer(buf))

	client := &http.Client{}
	request, error := http.NewRequest("POST", AUTH_SERVICE_URL+"/auth", rdr)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		http.Redirect(w, r, "/resource", http.StatusFound)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	authURL := map[string]interface{}{"authEndpoint": AUTH_SERVICE_URL + "/auth"}
	outputHTML(w, "../public/login.html", authURL)
}

func Resource(w http.ResponseWriter, r *http.Request) {
	outputHTML(w, "../public/resource.html", nil)
}

func outputHTML(w http.ResponseWriter, filename string, data interface{}) {
	t, err := template.ParseFiles(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func main() {
	http.HandleFunc("/", Login)
	http.HandleFunc("/authenticate", Authenticate)
	http.HandleFunc("/resource", Resource)
	log.Println("Serving at port 1234")
	http.ListenAndServe(":1234", nil)

}
