package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	RESOURCE_URL = "http://localhost:5555/get_resource"
)

func main() {
	fmt.Println("HTTP JSON POST URL:", RESOURCE_URL)

	var jsonData = []byte(`{
		"token": "12345"
	}`)
	request, error := http.NewRequest("POST", RESOURCE_URL, bytes.NewBuffer(jsonData))
	//request.Header.Set("Content-Type", "application/json; charset=UTF-8")

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
