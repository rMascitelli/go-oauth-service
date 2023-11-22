package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Handler struct {
	postgres PostgresConnector
}

func NewHandler(postgres PostgresConnector) Handler {
	return Handler{
		postgres: postgres,
	}
}

func (h *Handler) HandleUserLogin(req *http.Request) (error, LoginResponse) {
	uc := UserCredentialForm{}
	loginResponse := LoginResponse{}
	err := json.NewDecoder(req.Body).Decode(&uc)
	if err != nil {
		return errors.New("Failed to decode JSON body"), loginResponse
	}
	log.Printf("Got auth request: %+v\n", uc)

	// Get SHA256 string of user and pass
	// Make entry into DB
	err, queried_user := h.postgres.QueryUser(uc.Email, uc.Password)
	if err != nil {
		return errors.New("Failed to query user"), loginResponse
	}

	err, token := h.postgres.CreateAndStoreSessionToken(queried_user.userid)
	if err != nil {
		return fmt.Errorf("Error while creating session token, err: %v", err), loginResponse
	}
	loginResponse.Token = token
	return nil, loginResponse
}

func (h *Handler) HandleRegistry(req *http.Request) error {
	log.Printf("Got request for - %s\n", req.URL)
	qKey, qVal, _ := strings.Cut(req.URL.RawQuery, "=")
	if qKey == "registry_type" {
		if qVal == "user" {
			// RegisterUser()
			uc := UserCredentialForm{}
			err := json.NewDecoder(req.Body).Decode(&uc)
			if err != nil {
				return err
			}
			log.Printf("  Got User register request: %+v\n", uc)
			h.postgres.RegisterUser(uc.Email, uc.Password)
			log.Printf("  Registered user with email %s\n", uc.Email)
			return nil
		} else if qVal == "service" {
			// RegisterService()
			return errors.New("I dont know how to register services yet!")
		} else {
			return errors.New(fmt.Sprintf("Invalid registry_type: %s", qVal))
		}
	} else {
		errStr := fmt.Sprintf("  Unknown query params: %s\n", req.URL.RawQuery)
		log.Println("  " + errStr)
		return errors.New(errStr)
	}
}

func (h *Handler) HandleIntrospect(req *http.Request) (error, IntrospectResponse) {
	log.Println("Got introspect request, method: ", req.Method)
	introspectRequest := IntrospectRequest{}
	introspectResponse := IntrospectResponse{
		Active: false,
	}

	err := json.NewDecoder(req.Body).Decode(&introspectRequest)
	if err != nil {
		return errors.New("Failed to decode JSON body"), introspectResponse
	}

	if err = h.postgres.GetToken(introspectRequest.Token); err != nil {
		return errors.New(fmt.Sprintf("Failed to get token, err: %s", err.Error())), introspectResponse
	}

	log.Println("Introspect success!")
	introspectResponse.Active = true
	return nil, introspectResponse
}
