package main

import (
	"net/http"
	"fmt"
	"os"
	"log"
)

type Router struct {
	port int
}

func NewRouter(port int) Router {
	if port <= 0 {
		log.Fatalf("Cannot create server at port %d\n", port)
	}

	return Router {
		port: port,
	}
}

func (rt *Router) StartRouter() {
	fmt.Printf("Serving at port %d...\n", rt.port)
	// http.HandleFunc("/", HomePage)
	// http.HandleFunc("/register", rt.RegisterPage)
	// http.HandleFunc("/login", LoginPage)
	http.ListenAndServe(fmt.Sprintf(":%d", rt.port), nil)
}

func (rt *Router) HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Home page hit!")
	rt.outputHTML(w, r, "public/index.html")
}


func (rt *Router) RegisterPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Register page hit!")
	rt.outputHTML(w, r, "public/register.html")
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