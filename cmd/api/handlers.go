package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type Handler struct {
	postgres PostgresConnector
}

func NewHandler(postgres PostgresConnector) Handler {
	return Handler{
		postgres: postgres,
	}
}

func (h *Handler) HandleUserLogin(req *http.Request) (error, string) {
	uc := UserCredentialForm{}
	err := json.NewDecoder(req.Body).Decode(&uc)
	if err != nil {
		return errors.New("Failed to decode JSON body"), ""
	}
	log.Printf("Got auth request: %+v\n", uc)

	// Get SHA256 string of user and pass
	// Make entry into DB
	err, queried_user := h.postgres.QueryUser(uc.Email, uc.Password)
	if err != nil {
		return errors.New("Failed to query user"), ""
	}

	err, token := h.postgres.CreateAndStoreSessionToken(queried_user.userid)
	if err != nil {
		return errors.New("Error while creating session token"), ""
	}
	return nil, token
}
