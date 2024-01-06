package routes

import (
	"book-and-rate/pkg/handlers"
	middleware "book-and-rate/pkg/middlewares"
	"github.com/gorilla/mux"
)

func RateRoutes(router *mux.Router) {
	subRouter := router.PathPrefix("/rates").Subrouter()
	subRouter.Use(middleware.AuthenticationMiddleware)
	subRouter.HandleFunc("", handlers.CreateRateHandler).Methods("POST")
	subRouter.HandleFunc("/{id}", handlers.GetRateHandler).Methods("GET")
	subRouter.HandleFunc("/{id}", handlers.UpdateRateHandler).Methods("PUT")
	subRouter.HandleFunc("/{id}", handlers.DeleteRateHandler).Methods("DELETE")
	subRouter.HandleFunc("/restaurants/{restaurantId}/rates", handlers.GetRatesForRestaurant).Methods("GET")
	subRouter.HandleFunc("/restaurants/{restaurantId}/average-rating", handlers.GetAverageRatingForRestaurant).Methods("GET")
	subRouter.HandleFunc("/recent", handlers.GetRecentRatings).Queries("limit", "{limit}").Methods("GET")
}
