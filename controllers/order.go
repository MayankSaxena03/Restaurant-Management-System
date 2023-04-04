package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/MayankSaxena03/Restaurant-Management-System/database"
	"github.com/MayankSaxena03/Restaurant-Management-System/models"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var orderCollection = database.OpenCollection(database.Client, "order")
var tableCollection = database.OpenCollection(database.Client, "table")

func GetOrders(w http.ResponseWriter, r *http.Request) {
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

	var orders []models.Order
	cursor, err := orderCollection.Find(r.Context(), bson.M{}, options.Find().SetLimit(int64(limit)).SetSkip(int64(skip)))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}
	if err = cursor.All(r.Context(), &orders); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

func GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if !primitive.IsValidObjectID(id) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid ID")
		return
	}

	orderID, _ := primitive.ObjectIDFromHex(id)
	var order models.Order
	query := bson.M{
		"_id": orderID,
	}
	err := orderCollection.FindOne(r.Context(), query).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode("Order not found")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder models.Order
	err := json.NewDecoder(r.Body).Decode(&newOrder)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid request body")
		return
	}

	validate := validator.New()
	err = validate.Struct(newOrder)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	query := bson.M{
		"_id": newOrder.TableID,
	}

	count, err := tableCollection.CountDocuments(r.Context(), query)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}
	if count == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Table not found")
		return
	}

	newOrder.CreatedOn = time.Now()
	newOrder.UpdatedOn = time.Now()
	result, err := orderCollection.InsertOne(r.Context(), newOrder)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Inserted ID: " + result.InsertedID.(primitive.ObjectID).Hex())
}

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if !primitive.IsValidObjectID(id) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid ID")
		return
	}

	orderID, _ := primitive.ObjectIDFromHex(id)

	var UpdateOrder models.Order
	err := json.NewDecoder(r.Body).Decode(&UpdateOrder)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid request body")
		return
	}

	validate := validator.New()
	err = validate.Struct(UpdateOrder)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	tableId := UpdateOrder.TableID
	query := bson.M{
		"_id": tableId,
	}
	count, err := tableCollection.CountDocuments(r.Context(), query)
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
	}

	if count == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Table not found")
		return
	}

	query = bson.M{
		"_id": orderID,
	}

	update := bson.M{
		"$set": bson.M{
			"table_id":  UpdateOrder.TableID,
			"updatedOn": time.Now(),
		},
	}

	result, err := orderCollection.UpdateOne(r.Context(), query, update)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}
	if result.MatchedCount == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Order not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Order updated successfully")
}

func OrderItemOrderCreator(ctx context.Context, order models.Order) (primitive.ObjectID, error) {
	order.CreatedOn = time.Now()
	order.UpdatedOn = time.Now()

	result, err := orderCollection.InsertOne(context.Background(), order)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}
