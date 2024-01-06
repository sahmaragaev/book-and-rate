package handlers

import (
	"book-and-rate/pkg/db"
	"book-and-rate/pkg/models"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateRateHandler handles the creation of a new rate
func CreateRateHandler(w http.ResponseWriter, r *http.Request) {
    var rate models.Rate
    if err := json.NewDecoder(r.Body).Decode(&rate); err != nil {
        log.Printf("CreateRateHandler: Error decoding rate: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    rate.Date = time.Now() // Setting the rate date to current time
    result, err := db.RateCollection.InsertOne(context.Background(), rate)
    if err != nil {
        log.Printf("CreateRateHandler: Error inserting rate: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("CreateRateHandler: Rate created, ID: %v", result.InsertedID)
    json.NewEncoder(w).Encode(result)
}

// GetRateHandler retrieves a rate by ID
func GetRateHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    rateId, err := primitive.ObjectIDFromHex(params["id"])
    if err != nil {
        log.Printf("GetRateHandler: Error parsing ID: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var rate models.Rate
    if err := db.RateCollection.FindOne(context.Background(), bson.M{"_id": rateId}).Decode(&rate); err != nil {
        log.Printf("GetRateHandler: Error finding rate: %v", err)
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    log.Printf("GetRateHandler: Rate retrieved, ID: %v", rateId)
    json.NewEncoder(w).Encode(rate)
}

// UpdateRateHandler updates a rate's details
func UpdateRateHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    rateId, err := primitive.ObjectIDFromHex(params["id"])
    if err != nil {
        log.Printf("UpdateRateHandler: Error parsing ID: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var rate models.Rate
    if err := json.NewDecoder(r.Body).Decode(&rate); err != nil {
        log.Printf("UpdateRateHandler: Error decoding rate: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    _, err = db.RateCollection.UpdateOne(context.Background(), bson.M{"_id": rateId}, bson.M{"$set": rate})
    if err != nil {
        log.Printf("UpdateRateHandler: Error updating rate: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("UpdateRateHandler: Rate updated, ID: %v", rateId)
    w.WriteHeader(http.StatusNoContent)
}

// DeleteRateHandler deletes a rate by ID
func DeleteRateHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    rateId, err := primitive.ObjectIDFromHex(params["id"])
    if err != nil {
        log.Printf("DeleteRateHandler: Error parsing ID: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    _, err = db.RateCollection.DeleteOne(context.Background(), bson.M{"_id": rateId})
    if err != nil {
        log.Printf("DeleteRateHandler: Error deleting rate: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("DeleteRateHandler: Rate deleted, ID: %v", rateId)
    w.WriteHeader(http.StatusNoContent)
}

// GetRatesForRestaurant retrieves all rates for a specific restaurant
func GetRatesForRestaurant(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    restaurantId, err := primitive.ObjectIDFromHex(params["restaurantId"])
    if err != nil {
        log.Printf("GetRatesForRestaurant: Error parsing restaurant ID: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var rates []models.Rate
    cursor, err := db.RateCollection.Find(context.Background(), bson.M{"restaurantId": restaurantId})
    if err != nil {
        log.Printf("GetRatesForRestaurant: Error finding rates: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer cursor.Close(context.Background())

    for cursor.Next(context.Background()) {
        var rate models.Rate
        if err := cursor.Decode(&rate); err != nil {
            log.Printf("GetRatesForRestaurant: Error decoding rate: %v", err)
            continue
        }
        rates = append(rates, rate)
    }

    if err := cursor.Err(); err != nil {
        log.Printf("GetRatesForRestaurant: Cursor error: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("GetRatesForRestaurant: Successfully retrieved rates")
    json.NewEncoder(w).Encode(rates)
}

// GetAverageRatingForRestaurant calculates the average rating for a specific restaurant
func GetAverageRatingForRestaurant(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    restaurantId, err := primitive.ObjectIDFromHex(params["restaurantId"])
    if err != nil {
        log.Printf("GetAverageRatingForRestaurant: Error parsing restaurant ID: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    pipeline := []bson.M{
        {"$match": bson.M{"restaurantId": restaurantId}},
        {"$group": bson.M{"_id": "$restaurantId", "averageRating": bson.M{"$avg": "$rating"}}},
    }

    cursor, err := db.RateCollection.Aggregate(context.Background(), pipeline)
    if err != nil {
        log.Printf("GetAverageRatingForRestaurant: Error aggregating ratings: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer cursor.Close(context.Background())

    var results []bson.M
    if cursor.All(context.Background(), &results) != nil || len(results) == 0 {
        log.Printf("GetAverageRatingForRestaurant: No ratings found")
        json.NewEncoder(w).Encode(bson.M{"averageRating": 0})
        return
    }

    log.Printf("GetAverageRatingForRestaurant: Successfully calculated average rating")
    json.NewEncoder(w).Encode(results[0])
}

// GetRecentRatings retrieves the most recent ratings, limited to a specified number
func GetRecentRatings(w http.ResponseWriter, r *http.Request) {
    limitQuery := r.URL.Query().Get("limit")
    limit, err := strconv.Atoi(limitQuery)
    if err != nil || limit <= 0 {
        limit = 10 // Default to 10 if no valid limit is provided
    }

    var ratings []models.Rate
    cursor, err := db.RateCollection.Find(
        context.Background(),
        bson.M{},
        options.Find().SetSort(bson.M{"date": -1}).SetLimit(int64(limit)),
    )
    if err != nil {
        log.Printf("GetRecentRatings: Error finding recent ratings: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer cursor.Close(context.Background())

    for cursor.Next(context.Background()) {
        var rate models.Rate
        if err := cursor.Decode(&rate); err != nil {
            log.Printf("GetRecentRatings: Error decoding rate: %v", err)
            continue
        }
        ratings = append(ratings, rate)
    }

    if err := cursor.Err(); err != nil {
        log.Printf("GetRecentRatings: Cursor error: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("GetRecentRatings: Successfully retrieved recent ratings")
    json.NewEncoder(w).Encode(ratings)
}
