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
	pgc PostgresConnector
}

type UserCredentialForm struct {
	email	string
	password string	
}

func NewRouter(port int, pgc PostgresConnector) Router {
	if port <= 0 {
		log.Fatalf("Cannot create server at port %d\n", port)
	}

	return Router {
		port: port,
		pgc: pgc,
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
	http.ListenAndServe(fmt.Sprintf(":%d", rt.port), nil)
}

func (rt *Router) HomePage(w http.ResponseWriter, r *http.Request) {
	rt.outputHTML(w, r, "public/index.html")
}

func (rt *Router) LoginPage(w http.ResponseWriter, r *http.Request) {
	rt.outputHTML(w, r, "public/login.html")
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
	    err, _ := rt.pgc.QueryUser(email_hash, pass_hash)
	    if err != nil {
	    	log.Println("Failed to authorize, err: ", err) 
	    	http.Redirect(w, r, "/", http.StatusFound)
	    }
	    // Create access token using uc, store in session_token table
		http.Redirect(w, r, "/welcome.html", http.StatusFound)
	}
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
	    rt.pgc.RegisterUser(email_hash, pass_hash)
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