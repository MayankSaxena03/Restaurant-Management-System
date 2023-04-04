package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/MayankSaxena03/Restaurant-Management-System/models"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMenus(w http.ResponseWriter, r *http.Request) {
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
	query := bson.M{}
	cursor, err := menuCollection.Find(r.Context(), query, options.Find().SetLimit(int64(limit)).SetSkip(int64(skip)))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	var menus []models.Menu
	if err = cursor.All(r.Context(), &menus); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(menus)
}

func GetMenu(w http.ResponseWriter, r *http.Request) {
	// Get the menu ID from the URL
	vars := mux.Vars(r)
	id := vars["id"]
	if !primitive.IsValidObjectID(id) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid Menu ID")
		return
	}

	// Get the menu from the database
	var menu models.Menu
	menuID, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{
		"_id": menuID,
	}
	err := menuCollection.FindOne(r.Context(), query).Decode(&menu)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode("Menu not found")
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	// Return the menu
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(menu)
}

func CreateMenu(w http.ResponseWriter, r *http.Request) {
	var menu models.Menu
	err := json.NewDecoder(r.Body).Decode(&menu)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid request payload")
		return
	}

	// Validate the menu
	validate := validator.New()
	err = validate.Struct(menu)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	menu.CreatedOn = time.Now()
	menu.UpdatedOn = time.Now()

	// Insert the menu into the database
	result, err := menuCollection.InsertOne(r.Context(), menu)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	// Return the menu ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Inserted menu with ID: " + result.InsertedID.(primitive.ObjectID).Hex())
}

func UpdateMenu(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if !primitive.IsValidObjectID(id) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid Menu ID")
		return
	}

	var menu models.Menu
	err := json.NewDecoder(r.Body).Decode(&menu)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid request payload")
		return
	}

	// Validate the menu
	validate := validator.New()
	err = validate.Struct(menu)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	menuID, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{
		"_id": menuID,
	}
	update := bson.M{
		"$set": bson.M{
			"name":      menu.Name,
			"category":  menu.Category,
			"startDate": menu.StartDate,
			"endDate":   menu.EndDate,
			"updatedOn": time.Now(),
		},
	}
	result, err := menuCollection.UpdateOne(r.Context(), query, update)
	if err != nil {
		if result.MatchedCount == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode("Menu with given ID not found")
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode("Something went wrong")
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Updated menu with ID: " + id)
}
