package main

import (
	"book-and-rate/pkg/config"
	"book-and-rate/pkg/db"
	"book-and-rate/pkg/routes"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig("config/config.json")
	db.Connect(cfg.MongoDbUrl)
	db.InitializeCollections()

	router := mux.NewRouter()
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	routes.UserRoutes(router)
	routes.RestaurantRoutes(router)
	routes.BookingRoutes(router)
	routes.RateRoutes(router)
	routes.RefreshTokenRoutes(router)

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
