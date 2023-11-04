package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
)

//OAuth Specification is described in these RFC Articles:
//	https://www.rfc-editor.org/rfc/rfc6749
//	https://www.rfc-editor.org/rfc/rfc7662

type Router struct {
	port     int
	postgres PostgresConnector
}

type UserCredentialForm struct {
	Email    string
	Password string
}

func NewRouter(port int, postgres PostgresConnector) Router {
	if port <= 0 {
		log.Fatalf("Cannot create server at port %d\n", port)
	}

	return Router{
		port:     port,
		postgres: postgres,
	}
}

func (r *Router) StartRouter() {
	log.Printf("Serving at port %d...\n", r.port)
	http.HandleFunc("/register", r.RegisterUser)
	http.HandleFunc("/auth", r.Auth)
	http.HandleFunc("/introspect", r.Introspect)
	http.ListenAndServe(fmt.Sprintf(":%d", r.port), nil)
}

func (r *Router) Auth(w http.ResponseWriter, req *http.Request) {
	log.Println("Got auth request")
	if req.Method == "POST" {
		req.ParseForm()
		creds := UserCredentialForm{
			Email:    req.FormValue("email"),
			Password: req.FormValue("password"),
		}

		// Get SHA256 string of user and pass
		// Make entry into DB
		email_hash := hex.EncodeToString(getSHA256Hash(creds.Email))
		pass_hash := hex.EncodeToString(getSHA256Hash(creds.Password))
		err, uc := r.postgres.QueryUser(email_hash, pass_hash)
		if err != nil {
			log.Println("Failed to authorize, err: ", err)
			writeJSONResponse(w, 400, "Failure")
			return
		}

		err, _ = r.postgres.CreateAndStoreSessionToken(uc.userid)
		if err != nil {
			log.Println("Failed to create session token, err: ", err)
			writeJSONResponse(w, 400, "Failure")
			return
		}

		writeJSONResponse(w, 200, "Success!")
		return
	}
}

func (r *Router) Introspect(w http.ResponseWriter, req *http.Request) {
	log.Println("Got introspect request, method: ", req.Method)
	authRequest := struct {
		Token string `json:"token"`
	}{}
	introspectResponse := struct {
		ActiveStatus bool
	}{}
	err := json.NewDecoder(req.Body).Decode(&authRequest)
	if err != nil {
		introspectResponse.ActiveStatus = false
		writeJSONResponse(w, 400, introspectResponse)
		return
	}

	if err := r.postgres.GetToken(authRequest.Token); err != nil {
		log.Println("Failed to get token, err: ", err)
		writeJSONResponse(w, 400, introspectResponse)
		return
	}
	log.Println("Introspect success!")
	introspectResponse.ActiveStatus = true
	writeJSONResponse(w, 200, introspectResponse)
}

func (r *Router) RegisterUser(w http.ResponseWriter, req *http.Request) {
	log.Printf("Got register request\n")
	if req.Method == "POST" {
		var uc UserCredentialForm
		err := json.NewDecoder(req.Body).Decode(&uc)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			writeJSONResponse(w, 400, "Failure")
			return
		}
		log.Printf("Got register request: %+v\n", uc)

		// Get SHA256 string of user and pass
		// Make entry into DB
		email_hash := hex.EncodeToString(getSHA256Hash(uc.Email))
		pass_hash := hex.EncodeToString(getSHA256Hash(uc.Password))
		r.postgres.RegisterUser(email_hash, pass_hash)
		log.Printf("Registered user with email %s\n", uc.Email)
		writeJSONResponse(w, 200, "Success!")
	}
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, data any) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	var jsonResp []byte
	var err error
	if reflect.ValueOf(data).Kind() == reflect.String {
		resp := make(map[string]string)
		resp["message"] = data.(string)
		jsonResp, err = json.Marshal(resp)
	} else {
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
