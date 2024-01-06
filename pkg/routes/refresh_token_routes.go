package routes

import (
	"book-and-rate/pkg/handlers"
	"github.com/gorilla/mux"
)

func RefreshTokenRoutes(router *mux.Router) {
	router.HandleFunc("/refresh-token", handlers.RefreshTokenHandler).Methods("POST")
}
