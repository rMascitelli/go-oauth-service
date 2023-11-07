package main

import (
	"crypto/sha256"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
)

func writeJSONResponse(w http.ResponseWriter, statusCode int, data any) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	var jsonResp []byte
	var err error
	if reflect.ValueOf(data).Kind() == reflect.String {
		log.Println("Sending message - ", data)
		resp := make(map[string]string)
		resp["message"] = data.(string)
		jsonResp, err = json.Marshal(resp)
	} else {
		log.Printf("Sending message - %+v", data)
		jsonResp, err = json.Marshal(data)
	}
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	_, err = w.Write(jsonResp)
	if err != nil {
		log.Fatalf("Error happened when writing Json Response. Err: %s", err)
	}
}

func getSHA256Hash(s string) []byte {
	h := sha256.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return bs
}
