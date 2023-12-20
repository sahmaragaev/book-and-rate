package handlers

import (
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

// CreateUserHandler handles the creation of a new user
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("CreateUserHandler: Error decoding user data: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Printf("CreateUserHandler: Error hashing password: %v", err)
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

	result, err := db.UserCollection.InsertOne(context.Background(), user)
	if err != nil {
		log.Printf("CreateUserHandler: Error inserting new user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("CreateUserHandler: User created successfully: %v", result.InsertedID)
	json.NewEncoder(w).Encode(result)
}

// GetUserHandler retrieves a user by ID
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		log.Printf("GetUserHandler: Error parsing ID: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	if err := db.UserCollection.FindOne(context.Background(), bson.M{"_id": userId}).Decode(&user); err != nil {
		log.Printf("GetUserHandler: Error finding user: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Printf("GetUserHandler: User retrieved: %v", userId)
	json.NewEncoder(w).Encode(user)
}

// UpdateUserHandler updates a user's details
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		log.Printf("UpdateUserHandler: Error parsing ID: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("UpdateUserHandler: Error decoding user: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Password != "" {
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			log.Printf("UpdateUserHandler: Error hashing password: %v", err)
			http.Error(w, "Error processing password", http.StatusInternalServerError)
			return
		}
		user.Password = hashedPassword
	}

	_, err = db.UserCollection.UpdateOne(context.Background(), bson.M{"_id": userId}, bson.M{"$set": user})
	if err != nil {
		log.Printf("UpdateUserHandler: Error updating user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("UpdateUserHandler: User updated successfully: %v", userId)
	w.WriteHeader(http.StatusNoContent)
}

// DeleteUserHandler deletes a user by ID
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		log.Printf("DeleteUserHandler: Error parsing ID: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.UserCollection.DeleteOne(context.Background(), bson.M{"_id": userId})
	if err != nil {
		log.Printf("DeleteUserHandler: Error deleting user: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("DeleteUserHandler: User deleted successfully: %v", userId)
	w.WriteHeader(http.StatusNoContent)
}

// LoginUserHandler handles the login process for a user
func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var loginDetails models.Login
	if err := json.NewDecoder(r.Body).Decode(&loginDetails); err != nil {
		log.Printf("LoginUserHandler: Error decoding login details: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	err := db.UserCollection.FindOne(context.Background(), bson.M{"phone": loginDetails.Phone}).Decode(&user)
	if err != nil {
		log.Printf("LoginUserHandler: Error finding user: %v", err)
		http.Error(w, "Invalid phone number or password", http.StatusUnauthorized)
		return
	}

	err = utils.ComparePasswords(user.Password, loginDetails.Password)
	if err != nil {
		log.Printf("LoginUserHandler: Password does not match: %v", err)
		http.Error(w, "Invalid phone number or password", http.StatusUnauthorized)
		return
	}

	log.Printf("LoginUserHandler: User logged in successfully: %v", user.ID)
	// Return success response or JWT token as needed
	w.WriteHeader(http.StatusOK)
}
