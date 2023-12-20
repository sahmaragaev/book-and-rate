package routes

import (
    "book-and-rate/pkg/handlers"
    "github.com/gorilla/mux"
)

func BookingRoutes(router *mux.Router) {
    router.HandleFunc("/bookings", handlers.CreateBookingHandler).Methods("POST")
    router.HandleFunc("/bookings/{id}", handlers.GetBookingHandler).Methods("GET")
    router.HandleFunc("/bookings/{id}", handlers.UpdateBookingHandler).Methods("PUT")
    router.HandleFunc("/bookings/{id}", handlers.DeleteBookingHandler).Methods("DELETE")
    router.HandleFunc("/bookings/{id}/cancel", handlers.CancelBookingHandler).Methods("PUT")
    router.HandleFunc("/restaurants/{restaurantId}/bookings", handlers.GetBookingsForRestaurant).Methods("GET")
    router.HandleFunc("/restaurants/{restaurantId}/active-bookings", handlers.GetActiveBookingsForRestaurant).Methods("GET")
    router.HandleFunc("/users/{userId}/future-bookings", handlers.GetFutureBookingsForUser).Methods("GET")
    router.HandleFunc("/restaurants/{restaurantId}/past-bookings", handlers.GetPastBookingsForRestaurant).Methods("GET")
}