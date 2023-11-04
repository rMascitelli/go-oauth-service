package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const (
	RESOURCE_URL     = "http://localhost:5555"
	AUTH_SERVICE_URL = "http://localhost:8080"
)

type UserCredentialForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// TODO: What is a more secure way to store this access token?
var AccessToken string

func GetResource(w http.ResponseWriter, r *http.Request) {
	AccessToken = "123"
	fmt.Println("HTTP JSON POST URL:", RESOURCE_URL)

	var introspectReq = []byte(fmt.Sprintf(`{
		"token": "%s"
	}`, AccessToken))
	request, error := http.NewRequest("POST", RESOURCE_URL+"/access_resource", bytes.NewBuffer(introspectReq))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		http.Redirect(w, r, "/resource.html", http.StatusFound)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func SendRegisterRequest(w http.ResponseWriter, r *http.Request) {
	SendUserCredentialForm(w, r, AUTH_SERVICE_URL+"/register", "/", "/")
}

func SendAuthRequest(w http.ResponseWriter, r *http.Request) {
	SendUserCredentialForm(w, r, AUTH_SERVICE_URL+"/auth", "/welcome.html", "/")
}

func SendUserCredentialForm(w http.ResponseWriter, r *http.Request, endpointURL string, passRoute string, failRoute string) {
	log.Printf("Sending UserCred form request to %s...\n", endpointURL)
	r.ParseForm()
	creds := UserCredentialForm{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	client := &http.Client{}
	json, err := json.Marshal(creds)
	if err != nil {
		panic(err)
	}
	request, err := http.NewRequest("POST", endpointURL, bytes.NewBuffer(json))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		http.Redirect(w, r, passRoute, http.StatusFound)
	} else {
		http.Redirect(w, r, failRoute, http.StatusFound)
	}
}

// outputHTML meant for use with HTML Templates
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

//TODO:
//	Currently, client is sending AUTH requests directly to Auth Service and serving its own resource
//	Need to split interactions into 2 categories:
//		Auth service: 		Register, Login, and getting to Welcome.html
//		Resource service: 	Use token from Auth service to make request, get resource.html back

func main() {
	AccessToken = ""
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)
	http.HandleFunc("/send_register_request", SendRegisterRequest)
	http.HandleFunc("/send_auth_request", SendAuthRequest)
	http.HandleFunc("/get_resource", GetResource)
	log.Println("Serving at port 1234")
	http.ListenAndServe(":1234", nil)
}
