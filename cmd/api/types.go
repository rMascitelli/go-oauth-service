package main

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

type IntrospectRequest = LoginResponse

type LoginResponse struct {
	Token string `json:"token"`
}
