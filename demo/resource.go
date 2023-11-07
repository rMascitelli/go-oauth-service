package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	INTROSPECT_ENDPOINT = "http://localhost:5001/introspect"
)

type Token struct {
	Stringval string `json:"token"`
}

type IntrospectResponse struct {
	Active bool `json:"active"`
}

func AccessResource(w http.ResponseWriter, r *http.Request) {
	log.Println("Someone is trying to access the resource...")
	// Decode token
	var t Token
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		log.Println("Error decoding request, err: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Got token: %+v\n", t)

	// Make new request
	var jsonData = []byte(fmt.Sprintf(`{
		"token": "%s"
	}`, t.Stringval))
	request, error := http.NewRequest("POST", INTROSPECT_ENDPOINT, bytes.NewBuffer(jsonData))
	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		log.Println("Error running request, err: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer response.Body.Close()

	introspect_resp := IntrospectResponse{}
	log.Printf("Introspect response: %+v\n", introspect_resp)
	err = json.NewDecoder(response.Body).Decode(&introspect_resp)
	if err != nil {
		log.Println("Error decoding response, err: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if introspect_resp.Active {
		log.Println("Token is active, granting access")
		w.WriteHeader(http.StatusOK)
	} else {
		log.Println("Token is not active, denying access")
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Write(nil)
}

func main() {
	fmt.Println("Resource: Hello world!")

	http.HandleFunc("/access_resource", AccessResource)
	http.ListenAndServe(":5002", nil)
}
