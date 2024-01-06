package handlers

import (
	"book-and-rate/pkg/auth"
	"book-and-rate/pkg/config"
	"book-and-rate/pkg/db"
	"book-and-rate/pkg/models"
	"book-and-rate/pkg/utils"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateRestaurantHandler handles the creation of a new restaurant
func CreateRestaurantHandler(w http.ResponseWriter, r *http.Request) {
	var restaurant models.Restaurant
	if err := json.NewDecoder(r.Body).Decode(&restaurant); err != nil {
		log.Printf("CreateRestaurantHandler: Error decoding restaurant data: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(restaurant.Password)
	if err != nil {
		log.Printf("CreateRestaurantHandler: Error hashing password: %v", err)
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}
	restaurant.Password = hashedPassword

	result, err := db.RestaurantCollection.InsertOne(context.Background(), restaurant)
	if err != nil {
		log.Printf("CreateRestaurantHandler: Error inserting new restaurant: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("CreateRestaurantHandler: Restaurant created successfully: %v", result.InsertedID)
	json.NewEncoder(w).Encode(result)
}

// GetRestaurantsHandler lists all restaurants
func GetRestaurantsHandler(w http.ResponseWriter, r *http.Request) {
	var restaurants []models.Restaurant

	cursor, err := db.RestaurantCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Printf("GetRestaurantsHandler: Error finding restaurants: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var restaurant models.Restaurant
		if err := cursor.Decode(&restaurant); err != nil {
			log.Printf("GetRestaurantsHandler: Error decoding restaurant: %v", err)
			continue
		}
		restaurants = append(restaurants, restaurant)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("GetRestaurantsHandler: Cursor error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("GetRestaurantsHandler: Successfully retrieved restaurants")
	json.NewEncoder(w).Encode(restaurants)
}

// GetRestaurantHandler retrieves a restaurant by ID
func GetRestaurantHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	restaurantId, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		log.Printf("GetRestaurantHandler: Error parsing ID: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var restaurant models.Restaurant
	if err := db.RestaurantCollection.FindOne(context.Background(), bson.M{"_id": restaurantId}).Decode(&restaurant); err != nil {
		log.Printf("GetRestaurantHandler: Error finding restaurant: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Printf("GetRestaurantHandler: Restaurant retrieved: %v", restaurantId)
	json.NewEncoder(w).Encode(restaurant)
}

// UpdateRestaurantHandler updates a restaurant's details
func UpdateRestaurantHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	restaurantId, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		log.Printf("UpdateRestaurantHandler: Error parsing ID: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var restaurant models.Restaurant
	if err := json.NewDecoder(r.Body).Decode(&restaurant); err != nil {
		log.Printf("UpdateRestaurantHandler: Error decoding restaurant: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if restaurant.Password != "" {
		hashedPassword, err := utils.HashPassword(restaurant.Password)
		if err != nil {
			log.Printf("UpdateRestaurantHandler: Error hashing password: %v", err)
			http.Error(w, "Error processing password", http.StatusInternalServerError)
			return
		}
		restaurant.Password = hashedPassword
	}

	_, err = db.RestaurantCollection.UpdateOne(context.Background(), bson.M{"_id": restaurantId}, bson.M{"$set": restaurant})
	if err != nil {
		log.Printf("UpdateRestaurantHandler: Error updating restaurant: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("UpdateRestaurantHandler: Restaurant updated successfully: %v", restaurantId)
	w.WriteHeader(http.StatusNoContent)
}

// DeleteRestaurantHandler deletes a restaurant by ID
func DeleteRestaurantHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	restaurantId, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		log.Printf("DeleteRestaurantHandler: Error parsing ID: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.RestaurantCollection.DeleteOne(context.Background(), bson.M{"_id": restaurantId})
	if err != nil {
		log.Printf("DeleteRestaurantHandler: Error deleting restaurant: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("DeleteRestaurantHandler: Restaurant deleted successfully: %v", restaurantId)
	w.WriteHeader(http.StatusNoContent)
}

// LoginRestaurantHandler handles the login process for a restaurant
func LoginRestaurantHandler(w http.ResponseWriter, r *http.Request) {
	var loginDetails models.Login
	if err := json.NewDecoder(r.Body).Decode(&loginDetails); err != nil {
		log.Printf("LoginRestaurantHandler: Error decoding login details: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var restaurant models.Restaurant
	err := db.RestaurantCollection.FindOne(context.Background(), bson.M{"phone": loginDetails.Phone}).Decode(&restaurant)
	if err != nil {
		log.Printf("LoginRestaurantHandler: Error finding restaurant: %v", err)
		http.Error(w, "Invalid phone number or password", http.StatusUnauthorized)
		return
	}

	if err = utils.ComparePasswords(restaurant.Password, loginDetails.Password); err != nil {
		log.Printf("LoginRestaurantHandler: Password does not match: %v", err)
		http.Error(w, "Invalid phone number or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT Token
	cfg := config.LoadConfig("./config/config.json")
	accessToken, err := auth.GenerateToken(restaurant.ID.Hex(), *cfg)
	if err != nil {
		log.Printf("LoginRestaurantHandler: Error generating access token: %v", err)
		http.Error(w, "Error generating access token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(restaurant.ID.Hex(), *cfg)
	if err != nil {
		log.Printf("LoginRestaurantHandler: Error generating refresh token: %v", err)
		http.Error(w, "Error generating refresh token", http.StatusInternalServerError)
		return
	}

	log.Printf("LoginRestaurantHandler: Restaurant logged in successfully: %v", restaurant.ID)
	json.NewEncoder(w).Encode(map[string]string{"accessToken": accessToken, "refreshToken": refreshToken})
}
