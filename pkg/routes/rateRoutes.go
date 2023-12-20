package routes

import (
	"book-and-rate/pkg/handlers"
	"github.com/gorilla/mux"
)

func RateRoutes(router *mux.Router) {
	router.HandleFunc("/rates", handlers.CreateRateHandler).Methods("POST")
	router.HandleFunc("/rates/{id}", handlers.GetRateHandler).Methods("GET")
	router.HandleFunc("/rates/{id}", handlers.UpdateRateHandler).Methods("PUT")
	router.HandleFunc("/rates/{id}", handlers.DeleteRateHandler).Methods("DELETE")
	router.HandleFunc("/restaurants/{restaurantId}/rates", handlers.GetRatesForRestaurant).Methods("GET")
	router.HandleFunc("/restaurants/{restaurantId}/average-rating", handlers.GetAverageRatingForRestaurant).Methods("GET")
	router.HandleFunc("/rates/recent", handlers.GetRecentRatings).Queries("limit", "{limit}").Methods("GET")
}