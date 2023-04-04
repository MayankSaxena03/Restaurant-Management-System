package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/MayankSaxena03/Restaurant-Management-System/database"
	"github.com/MayankSaxena03/Restaurant-Management-System/helpers"
	"github.com/MayankSaxena03/Restaurant-Management-System/models"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var foodCollection = database.OpenCollection(database.Client, "food")
var menuCollection = database.OpenCollection(database.Client, "menu")

func GetFoods(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	l := queryParams.Get("limit")
	limit, err := strconv.Atoi(l)
	if err != nil || limit < 0 {
		limit = 10
	}
	s := queryParams.Get("skip")
	skip, err := strconv.Atoi(s)
	if err != nil || skip < 0 {
		skip = 0
	}
	var foods []models.Food
	cursor, err := foodCollection.Find(r.Context(), bson.M{}, options.Find().SetLimit(int64(limit)).SetSkip(int64(skip)))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}
	if err = cursor.All(r.Context(), &foods); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(foods)
}

func GetFood(w http.ResponseWriter, r *http.Request) {
	// Get the food ID from the URL
	vars := mux.Vars(r)
	id := vars["id"]
	if !primitive.IsValidObjectID(id) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid Food ID")
		return
	}

	// Get the food from the database
	var food models.Food
	foodID, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{
		"_id": foodID,
	}
	err := foodCollection.FindOne(r.Context(), query).Decode(&food)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode("Food not found")
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	// Return the food
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(food)
}

func CreateFood(w http.ResponseWriter, r *http.Request) {
	// Get the food from the request body
	var food models.Food
	err := json.NewDecoder(r.Body).Decode(&food)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid request body")
		return
	}

	// Validate the food
	validate := validator.New()
	err = validate.Struct(food)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	//Check if menu with ID exists
	query := bson.M{
		"_id": food.MenuID,
	}
	count, err := menuCollection.CountDocuments(r.Context(), query)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}
	if count == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Menu with ID does not exist")
		return
	}

	//Update Price Fields
	food.Price = helpers.ToFixed(food.Price, 2)

	//Set Time Fields
	food.CreatedOn = time.Now()
	food.UpdatedOn = time.Now()

	// Insert the food into the database
	result, err := foodCollection.InsertOne(r.Context(), food)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	// Return the food ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Inserted at ID : " + result.InsertedID.(primitive.ObjectID).Hex())
}

func UpdateFood(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if !primitive.IsValidObjectID(id) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid Food ID")
		return
	}

	// Get the food from the request body
	var food models.Food
	err := json.NewDecoder(r.Body).Decode(&food)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid request body")
		return
	}

	// Validate the food
	validate := validator.New()
	err = validate.Struct(food)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	//Check if menu with ID exists
	query := bson.M{
		"_id": food.MenuID,
	}
	count, err := menuCollection.CountDocuments(r.Context(), query)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}
	if count == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Menu with ID does not exist")
		return
	}

	//Update Price Fields
	food.Price = helpers.ToFixed(food.Price, 2)

	//Set Time Fields
	food.UpdatedOn = time.Now()

	// Update the food in the database
	foodID, _ := primitive.ObjectIDFromHex(id)
	query = bson.M{
		"_id": foodID,
	}
	update := bson.M{
		"$set": food,
	}
	result, err := foodCollection.UpdateOne(r.Context(), query, update)
	if err != nil {
		if result.MatchedCount == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode("Food not found")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	// Return the food ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Updated at ID : " + id)
}
