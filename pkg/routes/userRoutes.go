package routes

import (
	"book-and-rate/pkg/handlers"

	"github.com/gorilla/mux"
)

func UserRoutes(router *mux.Router) {
	router.HandleFunc("/users", handlers.CreateUserHandler).Methods("POST")
	router.HandleFunc("/users/login", handlers.LoginUserHandler).Methods("POST")
	router.HandleFunc("/users/{id}", handlers.GetUserHandler).Methods("GET")
	router.HandleFunc("/users/{id}", handlers.UpdateUserHandler).Methods("PUT")
	router.HandleFunc("/users/{id}", handlers.DeleteUserHandler).Methods("DELETE")
}
