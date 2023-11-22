package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"
)

const (
	AUTH_SERVICE_URL = "http://localhost:5001"
	RESULTS_LOGFILE  = "../../test_history.log"
)

// go clean -testcache; go test -v *.go

func IntrospectTest(token LoginResponse) error {
	var introspectResp IntrospectResponse
	endpointURL := AUTH_SERVICE_URL + "/introspect"
	client := &http.Client{}
	tokenJson, _ := json.Marshal(token)
	request, err := http.NewRequest("POST", endpointURL, bytes.NewBuffer(tokenJson))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	response, err := client.Do(request)
	if err != nil {
		return errors.New(fmt.Sprintf("Error making request, err: %v\n", err))
	}
	defer response.Body.Close()

	// Get token in response
	err = json.NewDecoder(response.Body).Decode(&introspectResp)
	if err != nil {
		return errors.New(fmt.Sprintf("Error decoding Token payload, err: %v\n", err))
	}

	// Validate payload
	if response.StatusCode != 200 || introspectResp.Active != true {
		return errors.New(fmt.Sprintf("Payload validation failed"))
	}
	return nil
}

func LoginTest(id int) (LoginResponse, error) {
	var retToken LoginResponse
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
		return retToken, errors.New(fmt.Sprintf("Error making request, err: %v\n", err))
	}
	defer response.Body.Close()

	// Get token in response
	err = json.NewDecoder(response.Body).Decode(&retToken)
	if err != nil {
		return retToken, errors.New(fmt.Sprintf("Error decoding Token payload, err: %v\n", err))
	}

	// Validate payload
	if response.StatusCode != 200 || retToken.Token == "" {
		return retToken, errors.New(fmt.Sprintf("Payload validation failed"))
	}
	return retToken, nil
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
	var wg sync.WaitGroup
	var i int
	k := 100
	numReqList := []int{1 * k}

	logFile, _ := os.OpenFile(RESULTS_LOGFILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	logFile.Write([]byte(fmt.Sprintf("[%s]\n", time.Now().Format(time.RFC822))))
	defer logFile.Close()

	for _, numReq := range numReqList {
		tokens := make([]LoginResponse, numReq+1)
		logFile.Write([]byte("---\n"))
		start := time.Now()
		for i = 0; i < numReq; i++ {
			wg.Add(1)
			go func() {
				err = RegisterTest(i)
				assert.Nil(t, err)
				wg.Done()
			}()
		}
		wg.Wait()
		elapsed_s := float64(time.Since(start).Milliseconds()) / 1000
		log := fmt.Sprintf("  Ran RegisterTest %d times in %.3f sec (%.2f req/s)\n", numReq, elapsed_s, float64(numReq)/elapsed_s)
		_, _ = logFile.Write([]byte(log))
		t.Logf(log)

		start = time.Now()
		for i = 0; i < numReq; i++ {
			wg.Add(1)
			go func() {
				token, err := LoginTest(i)
				assert.Nil(t, err)
				tokens[i] = token
				wg.Done()
			}()
		}
		wg.Wait()
		elapsed_s = float64(time.Since(start).Milliseconds()) / 1000
		log = fmt.Sprintf("  Ran LoginTest %d times in %.3f sec (%.2f req/s)\n", numReq, elapsed_s, float64(numReq)/elapsed_s)
		_, _ = logFile.Write([]byte(log))
		t.Logf(log)

		start = time.Now()
		for i = 0; i < numReq; i++ {
			wg.Add(1)
			go func() {
				err = IntrospectTest(tokens[i])
				assert.Nil(t, err)
				wg.Done()
			}()
		}
		wg.Wait()
		elapsed_s = float64(time.Since(start).Milliseconds()) / 1000
		log = fmt.Sprintf("  Ran IntrospectTest %d times in %.3f sec (%.2f req/s)\n", numReq, elapsed_s, float64(numReq)/elapsed_s)
		_, _ = logFile.Write([]byte(log))
		t.Logf(log)

	}

	logFile.Write([]byte("---\n"))
}
