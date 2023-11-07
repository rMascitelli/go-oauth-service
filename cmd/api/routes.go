package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

//OAuth Specification is described in these RFC Articles:
//	https://www.rfc-editor.org/rfc/rfc6749
//	https://www.rfc-editor.org/rfc/rfc7662

type Router struct {
	port    int
	handler Handler
}

type UserCredentialForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type IntrospectResponse struct {
	Active bool `json:"active"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

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
	err, token := r.handler.HandleUserLogin(req)
	if err != nil {
		log.Printf("Error occured while handling User Login, err: %v\n", err)
		writeJSONResponse(w, 400, err.Error())
	} else {
		log.Printf("Succesfully authenticated")
		loginResp := LoginResponse{
			Token: token,
		}
		writeJSONResponse(w, 200, loginResp)
	}
}

func (r *Router) Register(w http.ResponseWriter, req *http.Request) {
	log.Printf("Got request for - %s\n", req.URL)
	log.Printf("Got request for - %s\n", req.URL.RawQuery)
	qKey, qVal, _ := strings.Cut(req.URL.RawQuery, "=")
	if qKey == "registry_type" {
		if qVal == "user" {
			uc := UserCredentialForm{}
			err := json.NewDecoder(req.Body).Decode(&uc)
			if err != nil {
				log.Println("Failed to register user")
				writeJSONResponse(w, 400, "Failure")
				return
			}
			log.Printf("  Got User register request: %+v\n", uc)

			r.handler.postgres.RegisterUser(uc.Email, uc.Password)
			log.Printf("  Registered user with email %s\n", uc.Email)
			writeJSONResponse(w, 200, "Success!")
			url.ParseQuery("")
		} else if qVal == "service" {
			log.Println("  I dont know how to register services yet!")
		} else {
			log.Println("  Bad query param: ", qVal)
			writeJSONResponse(w, 400, "Failure!")
		}
	} else {
		errStr := fmt.Sprintf("  Unknown query params: %s\n", req.URL.RawQuery)
		log.Println(errStr)
		writeJSONResponse(w, 400, errStr)
	}
}

func (r *Router) Introspect(w http.ResponseWriter, req *http.Request) {
	log.Println("Got introspect request, method: ", req.Method)
	authRequest := struct {
		Token string `json:"token"`
	}{}
	introspectResponse := IntrospectResponse{}
	err := json.NewDecoder(req.Body).Decode(&authRequest)
	if err != nil {
		introspectResponse.Active = false
		writeJSONResponse(w, 400, introspectResponse)
		return
	}

	if err := r.handler.postgres.GetToken(authRequest.Token); err != nil {
		log.Println("Failed to get token, err: ", err)
		writeJSONResponse(w, 400, introspectResponse)
		return
	}
	log.Println("Introspect success!")
	introspectResponse.Active = true
	writeJSONResponse(w, 200, introspectResponse)
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, data any) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	var jsonResp []byte
	var err error
	if reflect.ValueOf(data).Kind() == reflect.String {
		log.Println("Sending message - ", data)
		resp := make(map[string]string)
		resp["message"] = data.(string)
		jsonResp, err = json.Marshal(resp)
	} else {
		log.Printf("Sending message - %+v", data)
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
