package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
)

const (
	AUTH_SERVICE_URL = "http://localhost:5001"
	RESULTS_LOGFILE  = "../../test_history.log"
)

// go clean -testcache; go test -v *.go

type Tester func(int) error

func LoginTest(id int) error {
	endpointURL := AUTH_SERVICE_URL + "/login"
	creds := UserCredentialForm{
		Email:    fmt.Sprintf("randuser-%d", id),
		Password: fmt.Sprintf("randpass-%d", id),
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

func RegisterTest(id int) error {
	endpointURL := AUTH_SERVICE_URL + "/register?registry_type=user"
	creds := UserCredentialForm{
		Email:    fmt.Sprintf("randuser-%d", id),
		Password: fmt.Sprintf("randpass-%d", id),
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

	// Validate payload
	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Payload validation failed"))
	}
	return nil
}

func TestStressor(t *testing.T) {
	var err error
	k := 1000
	numReq := []int{1 * k, 10 * k}
	funcList := []Tester{RegisterTest, LoginTest}
	logFile, _ := os.OpenFile(RESULTS_LOGFILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	logFile.Write([]byte(fmt.Sprintf("[%s]\n", time.Now().Format(time.RFC822))))
	defer logFile.Close()

	for _, maxReq := range numReq {
		logFile.Write([]byte("---\n"))
		for _, f := range funcList {
			_, funcName, _ := strings.Cut(runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), ".")
			start := time.Now()
			for i := 0; i < maxReq; i++ {
				err = f(i)
				assert.Nil(t, err)
			}
			elapsed_s := float64(time.Since(start).Milliseconds()) / 1000
			log := fmt.Sprintf("  Ran %s %d times in %.3f sec (%.2f req/s)\n", funcName, maxReq, elapsed_s, float64(maxReq)/elapsed_s)
			_, _ = logFile.Write([]byte(log))
			t.Logf(log)
		}
	}
	logFile.Write([]byte("---\n"))
}
