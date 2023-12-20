package handlers

import (
    "book-and-rate/pkg/db"
    "book-and-rate/pkg/models"
    "context"
    "encoding/json"
    "log"
    "net/http"
    "time"

    "github.com/gorilla/mux"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateBookingHandler handles the creation of a new booking
func CreateBookingHandler(w http.ResponseWriter, r *http.Request) {
    var booking models.Booking
    if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
        log.Printf("CreateBookingHandler: Error decoding booking: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    result, err := db.BookingCollection.InsertOne(context.Background(), booking)
    if err != nil {
        log.Printf("CreateBookingHandler: Error inserting booking: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("CreateBookingHandler: Booking created, ID: %v", result.InsertedID)
    json.NewEncoder(w).Encode(result)
}

// GetBookingHandler retrieves a booking by ID
func GetBookingHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    bookingId, err := primitive.ObjectIDFromHex(params["id"])
    if err != nil {
        log.Printf("GetBookingHandler: Error parsing ID: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var booking models.Booking
    if err := db.BookingCollection.FindOne(context.Background(), bson.M{"_id": bookingId}).Decode(&booking); err != nil {
        log.Printf("GetBookingHandler: Error finding booking: %v", err)
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    log.Printf("GetBookingHandler: Booking retrieved, ID: %v", bookingId)
    json.NewEncoder(w).Encode(booking)
}

// UpdateBookingHandler updates a booking's details
func UpdateBookingHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    bookingId, err := primitive.ObjectIDFromHex(params["id"])
    if err != nil {
        log.Printf("UpdateBookingHandler: Error parsing ID: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var booking models.Booking
    if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
        log.Printf("UpdateBookingHandler: Error decoding booking: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    _, err = db.BookingCollection.UpdateOne(context.Background(), bson.M{"_id": bookingId}, bson.M{"$set": booking})
    if err != nil {
        log.Printf("UpdateBookingHandler: Error updating booking: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("UpdateBookingHandler: Booking updated, ID: %v", bookingId)
    w.WriteHeader(http.StatusNoContent)
}

// DeleteBookingHandler deletes a booking by ID
func DeleteBookingHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    bookingId, err := primitive.ObjectIDFromHex(params["id"])
    if err != nil {
        log.Printf("DeleteBookingHandler: Error parsing ID: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    _, err = db.BookingCollection.DeleteOne(context.Background(), bson.M{"_id": bookingId})
    if err != nil {
        log.Printf("DeleteBookingHandler: Error deleting booking: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("DeleteBookingHandler: Booking deleted, ID: %v", bookingId)
    w.WriteHeader(http.StatusNoContent)
}

// CancelBookingHandler marks a booking as canceled (soft delete)
func CancelBookingHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    bookingId, err := primitive.ObjectIDFromHex(params["id"])
    if err != nil {
        log.Printf("CancelBookingHandler: Error parsing ID: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    update := bson.M{"$set": bson.M{"cancelled": true}}
    _, err = db.BookingCollection.UpdateOne(context.Background(), bson.M{"_id": bookingId}, update)
    if err != nil {
        log.Printf("CancelBookingHandler: Error canceling booking: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("CancelBookingHandler: Booking canceled, ID: %v", bookingId)
    w.WriteHeader(http.StatusNoContent)
}

// GetBookingsForRestaurant retrieves all bookings for a specific restaurant
func GetBookingsForRestaurant(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    restaurantId, err := primitive.ObjectIDFromHex(params["restaurantId"])
    if err != nil {
        log.Printf("GetBookingsForRestaurant: Error parsing restaurant ID: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var bookings []models.Booking
    cursor, err := db.BookingCollection.Find(context.Background(), bson.M{"restaurantId": restaurantId})
    if err != nil {
        log.Printf("GetBookingsForRestaurant: Error finding bookings: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer cursor.Close(context.Background())

    for cursor.Next(context.Background()) {
        var booking models.Booking
        if err := cursor.Decode(&booking); err != nil {
            log.Printf("GetBookingsForRestaurant: Error decoding booking: %v", err)
            continue
        }
        bookings = append(bookings, booking)
    }

    if err := cursor.Err(); err != nil {
        log.Printf("GetBookingsForRestaurant: Cursor error: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("GetBookingsForRestaurant: Successfully retrieved bookings")
    json.NewEncoder(w).Encode(bookings)
}

// GetBookingsForUser retrieves all bookings made by a specific user
func GetBookingsForUser(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    userId, err := primitive.ObjectIDFromHex(params["userId"])
    if err != nil {
        log.Printf("GetBookingsForUser: Error parsing user ID: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var bookings []models.Booking
    cursor, err := db.BookingCollection.Find(context.Background(), bson.M{"userId": userId})
    if err != nil {
        log.Printf("GetBookingsForUser: Error finding bookings: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer cursor.Close(context.Background())

    for cursor.Next(context.Background()) {
        var booking models.Booking
        if err := cursor.Decode(&booking); err != nil {
            log.Printf("GetBookingsForUser: Error decoding booking: %v", err)
            continue
        }
        bookings = append(bookings, booking)
    }

    if err := cursor.Err(); err != nil {
        log.Printf("GetBookingsForUser: Cursor error: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("GetBookingsForUser: Successfully retrieved bookings")
    json.NewEncoder(w).Encode(bookings)
}

// GetBookingsByDate retrieves bookings on a specific date
func GetBookingsByDate(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    date, err := time.Parse(time.RFC3339, params["date"])
    if err != nil {
        log.Printf("GetBookingsByDate: Error parsing date: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var bookings []models.Booking
    cursor, err := db.BookingCollection.Find(context.Background(), bson.M{"date": bson.M{"$gte": date, "$lt": date.AddDate(0, 0, 1)}})
    if err != nil {
        log.Printf("GetBookingsByDate: Error finding bookings: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer cursor.Close(context.Background())

    for cursor.Next(context.Background()) {
        var booking models.Booking
        if err := cursor.Decode(&booking); err != nil {
            log.Printf("GetBookingsByDate: Error decoding booking: %v", err)
            continue
        }
        bookings = append(bookings, booking)
    }

    if err := cursor.Err(); err != nil {
        log.Printf("GetBookingsByDate: Cursor error: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("GetBookingsByDate: Successfully retrieved bookings")
    json.NewEncoder(w).Encode(bookings)
}

// GetActiveBookingsForRestaurant retrieves active (non-canceled) bookings for a specific restaurant
func GetActiveBookingsForRestaurant(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    restaurantId, err := primitive.ObjectIDFromHex(params["restaurantId"])
    if err != nil {
        log.Printf("GetActiveBookingsForRestaurant: Error parsing restaurant ID: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var bookings []models.Booking
    cursor, err := db.BookingCollection.Find(context.Background(), bson.M{
        "restaurantId": restaurantId,
        "cancelled":    bson.M{"$ne": true},
    })
    if err != nil {
        log.Printf("GetActiveBookingsForRestaurant: Error finding bookings: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer cursor.Close(context.Background())

    for cursor.Next(context.Background()) {
        var booking models.Booking
        if err := cursor.Decode(&booking); err != nil {
            log.Printf("GetActiveBookingsForRestaurant: Error decoding booking: %v", err)
            continue
        }
        bookings = append(bookings, booking)
    }

    if err := cursor.Err(); err != nil {
        log.Printf("GetActiveBookingsForRestaurant: Cursor error: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("GetActiveBookingsForRestaurant: Successfully retrieved active bookings")
    json.NewEncoder(w).Encode(bookings)
}

// GetFutureBookingsForUser retrieves future bookings for a specific user
func GetFutureBookingsForUser(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    userId, err := primitive.ObjectIDFromHex(params["userId"])
    if err != nil {
        log.Printf("GetFutureBookingsForUser: Error parsing user ID: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var bookings []models.Booking
    cursor, err := db.BookingCollection.Find(context.Background(), bson.M{
        "userId": userId,
        "date":   bson.M{"$gte": time.Now()},
    })
    if err != nil {
        log.Printf("GetFutureBookingsForUser: Error finding bookings: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer cursor.Close(context.Background())

    for cursor.Next(context.Background()) {
        var booking models.Booking
        if err := cursor.Decode(&booking); err != nil {
            log.Printf("GetFutureBookingsForUser: Error decoding booking: %v", err)
            continue
        }
        bookings = append(bookings, booking)
    }

    if err := cursor.Err(); err != nil {
        log.Printf("GetFutureBookingsForUser: Cursor error: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("GetFutureBookingsForUser: Successfully retrieved future bookings")
    json.NewEncoder(w).Encode(bookings)
}

// GetPastBookingsForRestaurant retrieves past bookings for a specific restaurant
func GetPastBookingsForRestaurant(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    restaurantId, err := primitive.ObjectIDFromHex(params["restaurantId"])
    if err != nil {
        log.Printf("GetPastBookingsForRestaurant: Error parsing restaurant ID: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var bookings []models.Booking
    cursor, err := db.BookingCollection.Find(context.Background(), bson.M{
        "restaurantId": restaurantId,
        "date":         bson.M{"$lt": time.Now()},
    })
    if err != nil {
        log.Printf("GetPastBookingsForRestaurant: Error finding bookings: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer cursor.Close(context.Background())

    for cursor.Next(context.Background()) {
        var booking models.Booking
        if err := cursor.Decode(&booking); err != nil {
            log.Printf("GetPastBookingsForRestaurant: Error decoding booking: %v", err)
            continue
        }
        bookings = append(bookings, booking)
    }

    if err := cursor.Err(); err != nil {
        log.Printf("GetPastBookingsForRestaurant: Cursor error: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("GetPastBookingsForRestaurant: Successfully retrieved past bookings")
    json.NewEncoder(w).Encode(bookings)
}