package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	//"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var CurrentToken string

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

func (rt *Router) StartRouter() {
	log.Printf("Serving at port %d...\n", rt.port)
	http.HandleFunc("/register_user", rt.RegisterUser)
	http.HandleFunc("/auth", rt.Auth)
	http.HandleFunc("/introspect", rt.Introspect)
	http.ListenAndServe(fmt.Sprintf(":%d", rt.port), nil)
}

func (rt *Router) Auth(w http.ResponseWriter, r *http.Request) {
	log.Println("Got auth request")
	if r.Method == "POST" {
		r.ParseForm()
		creds := UserCredentialForm{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		// Get SHA256 string of user and pass
		// Make entry into DB
		email_hash := hex.EncodeToString(getSHA256Hash(creds.Email))
		pass_hash := hex.EncodeToString(getSHA256Hash(creds.Password))
		err, uc := rt.postgres.QueryUser(email_hash, pass_hash)
		if err != nil {
			log.Println("Failed to authorize, err: ", err)
			writeJSONResponse(w, 400, "Failure")
			return
		}

		err, _ = rt.postgres.CreateAndStoreSessionToken(uc.userid)
		if err != nil {
			log.Println("Failed to create session token, err: ", err)
			writeJSONResponse(w, 400, "Failure")
			return
		}

		writeJSONResponse(w, 200, "Success!")
		return
	}
}

func (rt *Router) Introspect(w http.ResponseWriter, r *http.Request) {
	log.Println("Got introspect request, method: ", r.Method)
	authRequest := struct {
		Token string `json:"token"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&authRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := rt.postgres.GetToken(authRequest.Token); err != nil {
		log.Println("Failed to get token, err: ", err)
		http.Redirect(w, r, "/", http.StatusFound)
	}
	//rt.outputHTML(w, r, "public/resource.html")
	http.Redirect(w, r, "/resource.html", http.StatusFound)
}

func (rt *Router) RegisterUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Got Register request")
	if r.Method == "POST" {
		r.ParseForm()
		creds := UserCredentialForm{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		// Get SHA256 string of user and pass
		// Make entry into DB
		email_hash := hex.EncodeToString(getSHA256Hash(creds.Email))
		pass_hash := hex.EncodeToString(getSHA256Hash(creds.Password))
		rt.postgres.RegisterUser(email_hash, pass_hash)
		log.Printf("Registered user with email %s\n", creds.Email)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, err := json.Marshal(resp)
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
