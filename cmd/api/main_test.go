package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

const (
	AUTH_SERVICE_URL = "http://localhost:5001"
)

func LoginTest() error {
	endpointURL := AUTH_SERVICE_URL + "/login"
	creds := UserCredentialForm{
		Email:    "root",
		Password: "dev",
	}
	client := &http.Client{}
	credsJson, _ := json.Marshal(creds)
	request, err := http.NewRequest("POST", endpointURL, bytes.NewBuffer(credsJson))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	response, err := client.Do(request)
	if err != nil {
		return errors.New(fmt.Sprintf("Error making request, err: %v\n", err))
	}
	defer response.Body.Close()

	// Get token in response
	tempToken := LoginResponse{}
	err = json.NewDecoder(response.Body).Decode(&tempToken)
	if err != nil {
		return errors.New(fmt.Sprintf("Error decoding Token payload, err: %v\n", err))
	}

	// Validate payload
	if response.StatusCode != 200 || tempToken.Token == "" {
		return errors.New(fmt.Sprintf("Payload validation failed"))
	}
	return nil
}

func TestStarter(t *testing.T) {
	go func() {
		postgres := NewPostgresConnector(true)
		rt := NewRouter(5001, postgres)
		rt.StartRouter()
	}()
	time.Sleep(time.Second * 3)

	var err error
	err = LoginTest()
	assert.Nil(t, err)
}
