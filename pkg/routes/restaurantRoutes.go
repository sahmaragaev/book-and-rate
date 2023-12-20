package routes

import (
	"book-and-rate/pkg/handlers"
	"github.com/gorilla/mux"
)

func RestaurantRoutes(router *mux.Router) {
	router.HandleFunc("/restaurants", handlers.CreateRestaurantHandler).Methods("POST")
	router.HandleFunc("/restaurants", handlers.GetRestaurantsHandler).Methods("GET")
	router.HandleFunc("/restaurants/{id}", handlers.GetRestaurantHandler).Methods("GET")
	router.HandleFunc("/restaurants/{id}", handlers.UpdateRestaurantHandler).Methods("PUT")
	router.HandleFunc("/restaurants/{id}", handlers.DeleteRestaurantHandler).Methods("DELETE")
	router.HandleFunc("/restaurants/login", handlers.LoginRestaurantHandler).Methods("POST")
}