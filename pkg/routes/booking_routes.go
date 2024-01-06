package routes

import (
	"book-and-rate/pkg/handlers"
	middleware "book-and-rate/pkg/middlewares"
	"github.com/gorilla/mux"
)

func BookingRoutes(router *mux.Router) {
	router.HandleFunc("/bookings", handlers.CreateBookingHandler).Methods("POST")

	subRouter := router.PathPrefix("/bookings").Subrouter()
	subRouter.Use(middleware.AuthenticationMiddleware)
	subRouter.HandleFunc("/{id}", handlers.GetBookingHandler).Methods("GET")
	subRouter.HandleFunc("/{id}", handlers.UpdateBookingHandler).Methods("PUT")
	subRouter.HandleFunc("/{id}", handlers.DeleteBookingHandler).Methods("DELETE")
	subRouter.HandleFunc("/{id}/cancel", handlers.CancelBookingHandler).Methods("PUT")
	subRouter.HandleFunc("/restaurants/{restaurantId}/bookings", handlers.GetBookingsForRestaurant).Methods("GET")
	subRouter.HandleFunc("/restaurants/{restaurantId}/active-bookings", handlers.GetActiveBookingsForRestaurant).Methods("GET")
	subRouter.HandleFunc("/users/{userId}/future-bookings", handlers.GetFutureBookingsForUser).Methods("GET")
	subRouter.HandleFunc("/restaurants/{restaurantId}/past-bookings", handlers.GetPastBookingsForRestaurant).Methods("GET")
}
