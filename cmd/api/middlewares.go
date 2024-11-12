package main

import (
	"net/http"

	"github.com/rs/cors"
)

func (app *application) enableCORS(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corsOptions := cors.New(cors.Options{
			AllowedOrigins:   app.config.cors.allowedOrigins,
			AllowedHeaders:   app.config.cors.allowedHeaders,
			AllowedMethods:   app.config.cors.allowedMethods,
			AllowCredentials: true,
		})

		next = corsOptions.Handler(next)

		next.ServeHTTP(w, r)
	})
}
