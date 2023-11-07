package main

import (
	"fmt"
	"log"
	"net/http"
)

//OAuth Specification is described in these RFC Articles:
//	https://www.rfc-editor.org/rfc/rfc6749
//	https://www.rfc-editor.org/rfc/rfc7662

func NewRouter(port int, postgres PostgresConnector) Router {
	if port <= 0 {
		log.Fatalf("Cannot create server at port %d\n", port)
	}

	return Router{
		port:    port,
		handler: NewHandler(postgres),
	}
}

func (r *Router) StartRouter() {
	log.Printf("Serving at port %d...\n", r.port)
	http.HandleFunc("/register", r.Register)
	http.HandleFunc("/login", r.Login)
	http.HandleFunc("/introspect", r.Introspect)
	http.ListenAndServe(fmt.Sprintf(":%d", r.port), nil)
}

func (r *Router) Login(w http.ResponseWriter, req *http.Request) {
	err, loginResponse := r.handler.HandleUserLogin(req)
	if err != nil {
		log.Printf("Error occured while handling User Login, err:\n  %v\n", err)
		writeJSONResponse(w, 400, err.Error())
	} else {
		log.Printf("Succesfully authenticated")
		writeJSONResponse(w, 200, loginResponse)
	}
}

func (r *Router) Register(w http.ResponseWriter, req *http.Request) {
	if err := r.handler.HandleRegistry(req); err != nil {
		writeJSONResponse(w, 400, err.Error())
	} else {
		writeJSONResponse(w, 200, "Success!")
	}
}

func (r *Router) Introspect(w http.ResponseWriter, req *http.Request) {
	err, introspectResponse := r.handler.HandleIntrospect(req)
	if err != nil {
		log.Printf("Error occured while handling Introspect, err:\n  %v\n", err)
		writeJSONResponse(w, 400, err.Error())
	} else {
		log.Printf("Succesfully authenticated")
		writeJSONResponse(w, 200, introspectResponse)
	}
}
