package core

import (
	"net/http"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes - list of route
type Routes []route

// MetaData of HTTP API response
type MetaData struct {
	Code      int    `json:"code"`
	RequestID string `json:"requestId"`
	Message   string `json:"message"`
}

//Response - complete structure of  HTTP response Meta + Data
type Response struct {
	Meta MetaData    `json:"meta"`
	Data interface{} `json:"data,omitempty"`
}
