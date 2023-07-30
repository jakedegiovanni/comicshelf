package main

import (
	"log"
	"net/http"
	"runtime/debug"
)

func RecoverHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(uri string) {
			if r := recover(); r != nil {
				log.Println("recovered handler", uri, r)
				debug.PrintStack()
			}
		}(r.URL.String())

		next.ServeHTTP(w, r)
	}
}
