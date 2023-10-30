package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	INTROSPECT_URL = "http://localhost:8080/introspect"
)

type Token struct {
	Stringval string `json:"token"`
}

func GetResource(w http.ResponseWriter, r *http.Request) {
	var t Token
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Got token: %+v\n", t)

	var jsonData = []byte(fmt.Sprintf(`{
		"token": "%s"
	}`, t.Stringval))
	request, error := http.NewRequest("POST", INTROSPECT_URL, bytes.NewBuffer(jsonData))
	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	introspect_resp := struct {
		Active string
	}{
		Active: "",
	}
	err = json.NewDecoder(response.Body).Decode(&introspect_resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Got introspect response: %+v\n", introspect_resp)
}

func main() {
	fmt.Println("Resource: Hello world!")

	http.HandleFunc("/get_resource", GetResource)
	http.ListenAndServe(":5555", nil)
}
