package routes

import (
	"book-and-rate/pkg/handlers"
	middleware "book-and-rate/pkg/middlewares"
	"github.com/gorilla/mux"
)

func RestaurantRoutes(router *mux.Router) {
	router.HandleFunc("/restaurants", handlers.CreateRestaurantHandler).Methods("POST")
	router.HandleFunc("/restaurants/login", handlers.LoginRestaurantHandler).Methods("POST")
	subRouter := router.PathPrefix("/restaurants").Subrouter()
	subRouter.Use(middleware.AuthenticationMiddleware)
	subRouter.HandleFunc("", handlers.GetRestaurantsHandler).Methods("GET")
	subRouter.HandleFunc("/{id}", handlers.GetRestaurantHandler).Methods("GET")
	subRouter.HandleFunc("/{id}", handlers.UpdateRestaurantHandler).Methods("PUT")
	subRouter.HandleFunc("/{id}", handlers.DeleteRestaurantHandler).Methods("DELETE")
}
