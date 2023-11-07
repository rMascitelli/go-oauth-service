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
	RESOURCE_URL     = "http://localhost:5002"
	AUTH_SERVICE_URL = "http://localhost:5001"
)

type UserCredentialForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

var activeToken TokenResponse

func GetResource(w http.ResponseWriter, req *http.Request) {
	fmt.Println("HTTP JSON POST URL:", RESOURCE_URL)
	tokenJson, err := json.Marshal(activeToken)
	if err != nil {
		log.Println("Failed to marshal activeToken")
	}
	request, error := http.NewRequest("POST", RESOURCE_URL+"/access_resource", bytes.NewBuffer(tokenJson))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		http.Redirect(w, req, "/resource.html", http.StatusFound)
	} else {
		http.Redirect(w, req, "/", http.StatusFound)
	}
}

func SendRegisterRequest(w http.ResponseWriter, r *http.Request) {
	endpointURL := AUTH_SERVICE_URL + "/register?registry_type=user"
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
	http.Redirect(w, r, "/", http.StatusFound)
}

func SendAuthRequest(w http.ResponseWriter, req *http.Request) {
	endpointURL := AUTH_SERVICE_URL + "/login"
	log.Printf("Sending UserCred form request to %s...\n", endpointURL)

	//  Parse login form
	req.ParseForm()
	creds := UserCredentialForm{
		Email:    req.FormValue("email"),
		Password: req.FormValue("password"),
	}

	client := &http.Client{}
	credsJson, _ := json.Marshal(creds)
	request, err := http.NewRequest("POST", endpointURL, bytes.NewBuffer(credsJson))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	// Get token in response
	tempToken := TokenResponse{}
	err = json.NewDecoder(response.Body).Decode(&tempToken)
	if err != nil {
		log.Printf("Error decoding Token payload, err: %v\n", err)
		return
	}

	activeToken.Token = tempToken.Token
	log.Println("Set active token to", activeToken.Token[:3])
	if response.StatusCode == 200 {
		http.Redirect(w, req, "/welcome.html", http.StatusFound)
	} else {
		http.Redirect(w, req, "/", http.StatusFound)
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

func main() {
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)
	http.HandleFunc("/send_register_request", SendRegisterRequest)
	http.HandleFunc("/send_auth_request", SendAuthRequest)
	http.HandleFunc("/get_resource", GetResource)
	log.Println("Serving at port 1234")
	http.ListenAndServe(":1234", nil)
}
