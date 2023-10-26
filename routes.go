package main

import (
	"net/http"
	"fmt"
	"os"
	"log"
	"crypto/sha256"
	"encoding/hex"
)

type Router struct {
	port int
	postgres PostgresConnector
}

type UserCredentialForm struct {
	email	string
	password string	
}

func NewRouter(port int, postgres PostgresConnector) Router {
	if port <= 0 {
		log.Fatalf("Cannot create server at port %d\n", port)
	}

	return Router {
		port: port,
		postgres: postgres,
	}
}

func (rt *Router) StartRouter() {
	log.Printf("Serving at port %d...\n", rt.port)
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.Handle("/login", fs)
	http.Handle("/register", fs)
	http.Handle("/welcome", fs)

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
	        email:   r.FormValue("email"),
	        password: r.FormValue("password"),
	    }

	    // Get SHA256 string of user and pass
	    // Make entry into DB
	    email_hash := hex.EncodeToString(getSHA256Hash(creds.email))
	    pass_hash := hex.EncodeToString(getSHA256Hash(creds.password))
	    err, uc := rt.postgres.QueryUser(email_hash, pass_hash)
	    if err != nil {
	    	log.Println("Failed to authorize, err: ", err) 
	    	http.Redirect(w, r, "/", http.StatusFound)
	    }
	    rt.postgres.CreateSessionToken(uc.userid)
		http.Redirect(w, r, "/welcome.html", http.StatusFound)
	}
}

func (rt *Router) Introspect(w http.ResponseWriter, r *http.Request) {
	log.Println("Got introspect request")
	// if r.Method == "POST" {

	// }
}

func (rt *Router) RegisterUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Got Register request")
	if r.Method == "POST" {
		r.ParseForm() 
		creds := UserCredentialForm{
	        email:   r.FormValue("email"),
	        password: r.FormValue("password"),
	    }

	    // Get SHA256 string of user and pass
	    // Make entry into DB
	    email_hash := hex.EncodeToString(getSHA256Hash(creds.email))
	    pass_hash := hex.EncodeToString(getSHA256Hash(creds.password))
	    rt.postgres.RegisterUser(email_hash, pass_hash)
	    log.Printf("Registered user with email %s\n", creds.email)
	    http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (rt *Router) outputHTML(w http.ResponseWriter, r *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, r, file.Name(), fi.ModTime(), file)
}

func getSHA256Hash(s string) []byte {
	h := sha256.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return bs
}