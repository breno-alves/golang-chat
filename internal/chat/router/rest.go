package router

import (
	"fmt"
	"log/slog"
	"net/http"
)

// Get wraps the router for GET method
func (router *Router) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	slog.Info(fmt.Sprintf("GET %s", path))
	router.Router.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for POST method
func (router *Router) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	slog.Info(fmt.Sprintf("POST %s", path))
	router.Router.HandleFunc(path, f).Methods("POST")
}

// Put wraps the router for PUT method
func (router *Router) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	slog.Info(fmt.Sprintf("PUT %s", path))
	router.Router.HandleFunc(path, f).Methods("PUT")
}

// Delete wraps the router for DELETE method
func (router *Router) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	slog.Info(fmt.Sprintf("DELETE %s", path))
	router.Router.HandleFunc(path, f).Methods("DELETE")
}
