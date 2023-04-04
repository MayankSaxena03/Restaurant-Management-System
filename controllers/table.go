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

func GetTables(w http.ResponseWriter, r *http.Request) {
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

	var tables []models.Table
	cursor, err := tableCollection.Find(r.Context(), bson.M{}, options.Find().SetLimit(int64(limit)).SetSkip(int64(skip)))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	if err = cursor.All(r.Context(), &tables); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tables)
}

func GetTable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tableID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid Table ID")
		return
	}

	var table models.Table
	err = tableCollection.FindOne(r.Context(), bson.M{"_id": tableID}).Decode(&table)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode("Table not found")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(table)
}

func CreateTable(w http.ResponseWriter, r *http.Request) {
	var table models.Table
	err := json.NewDecoder(r.Body).Decode(&table)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid Table")
		return
	}

	validate := validator.New()
	err = validate.Struct(table)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	table.CreatedOn = time.Now()
	table.UpdatedOn = time.Now()

	result, err := tableCollection.InsertOne(r.Context(), table)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Inserted with ID: " + result.InsertedID.(primitive.ObjectID).Hex())
}

func UpdateTable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tableID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid Table ID")
		return
	}

	var table models.Table
	err = json.NewDecoder(r.Body).Decode(&table)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid Table")
		return
	}

	validate := validator.New()
	err = validate.Struct(table)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	table.UpdatedOn = time.Now()

	query := bson.M{
		"_id": tableID,
	}
	update := bson.M{
		"$set": table,
	}

	result, err := tableCollection.UpdateOne(r.Context(), query, update)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	if result.ModifiedCount == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Table not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
}
