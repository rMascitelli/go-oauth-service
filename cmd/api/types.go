package main

import (
	proclient "github.com/rMascitelli/go-prometheus-metrics-helper/client"
)

const (
	NUM_REGISTRY          = "num_registry_requests"
	NUM_LOGIN             = "num_login_requests"
	NUM_INTROSPECT        = "num_introspect_requests"
	ELAPSED_REGISTRY_MS   = "time_on_registry_request_ms"
	ELAPSED_LOGIN_MS      = "time_on_login_request_ms"
	ELAPSED_INTROSPECT_MS = "time_on_introspect_request_ms"
	SERVICENAME           = "oauth-service"
)

type Router struct {
	port          int
	handler       Handler
	metricsClient proclient.PrometheusClient
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
