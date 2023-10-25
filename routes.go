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
	http.HandleFunc("/", rt.HomePage)
	http.HandleFunc("/register", rt.RegisterPage)
	// http.HandleFunc("/login", rt.LoginPage)
	http.ListenAndServe(fmt.Sprintf(":%d", rt.port), nil)
}

func (rt *Router) HomePage(w http.ResponseWriter, r *http.Request) {
	log.Println("Home page hit!")
	rt.outputHTML(w, r, "public/index.html")
}

func (rt *Router) RegisterPage(w http.ResponseWriter, r *http.Request) {
	rt.outputHTML(w, r, "public/register.html")
	r.ParseForm()
	log.Println("Register page hit!")   
	creds := UserCredentialForm{
        email:   r.FormValue("email"),
        password: r.FormValue("password"),
    }
    log.Printf("Rcvd registration:\n  Email: %s\n  Password: %s\n", creds.email, creds.password)

    // Get SHA256 string of user and pass
    // Make entry into DB
    email_hash := hex.EncodeToString(getSHA256Hash(creds.email))
    pass_hash := hex.EncodeToString(getSHA256Hash(creds.password))
    rt.pgc.RegisterUser(email_hash, pass_hash)
    rt.pgc.QueryUser(email_hash, pass_hash)
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