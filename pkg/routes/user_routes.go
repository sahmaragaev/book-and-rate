package routes

import (
	"book-and-rate/pkg/handlers"
	middleware "book-and-rate/pkg/middlewares"

	"github.com/gorilla/mux"
)

func UserRoutes(router *mux.Router) {
	router.HandleFunc("/users", handlers.CreateUserHandler).Methods("POST")
	router.HandleFunc("/users/login", handlers.LoginUserHandler).Methods("POST")
	subRouter := router.PathPrefix("/users").Subrouter()
	subRouter.Use(middleware.AuthenticationMiddleware)
	subRouter.HandleFunc("/{id}", handlers.GetUserHandler).Methods("GET")
	subRouter.HandleFunc("/{id}", handlers.UpdateUserHandler).Methods("PUT")
	subRouter.HandleFunc("/{id}", handlers.DeleteUserHandler).Methods("DELETE")
}
