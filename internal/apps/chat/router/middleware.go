package router

import (
	"context"
	"net/http"
)

func (router *Router) handleRestRequest(handler RequestHandlerFunction) http.HandlerFunc {
	ctx := context.Background()
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS,PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "OPTIONS" {
			return
		}
		handler(ctx, router.Db, router.Cache, w, r)
	}
}

func (router *Router) handleWsRequest(handler RequestHandlerFunction) http.HandlerFunc {
	ctx := context.Background()
	return func(w http.ResponseWriter, r *http.Request) {
		handler(ctx, router.Db, router.Cache, w, r)
	}
}
