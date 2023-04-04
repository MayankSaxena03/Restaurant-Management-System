package controllers

import (
	"context"
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

type OrderItemPack struct {
	TableID    primitive.ObjectID `json:"tableID,omitempty" bson:"tableID,omitempty"`
	OrderItems []models.OrderItem `json:"orderItems,omitempty" bson:"orderItems,omitempty"`
}

var orderItemCollection = database.OpenCollection(database.Client, "orderItem")

func GetOrderItems(w http.ResponseWriter, r *http.Request) {
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

	var orderItems []models.OrderItem
	cursor, err := orderItemCollection.Find(r.Context(), bson.M{}, options.Find().SetLimit(int64(limit)).SetSkip(int64(skip)))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}
	if err = cursor.All(r.Context(), &orderItems); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orderItems)
}

func GetOrderItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderItemID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid Order Item ID")
		return
	}

	var orderItem models.OrderItem
	query := bson.M{
		"_id": orderItemID,
	}

	if err = orderItemCollection.FindOne(r.Context(), query).Decode(&orderItem); err != nil {
		if err == mongo.ErrNoDocuments {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode("Order Item not found")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orderItem)
}

func GetOrderItemsByOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := primitive.ObjectIDFromHex(vars["orderId"])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid Order ID")
		return
	}

	allOrderItems, err := ItemsByOrder(r.Context(), orderID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allOrderItems)
}

func CreateOrderItem(w http.ResponseWriter, r *http.Request) {
	var orderItemPack OrderItemPack
	var order models.Order

	if err := json.NewDecoder(r.Body).Decode(&orderItemPack); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid Order Item")
		return
	}

	order.CreatedOn = time.Now()
	orderItemsToBeInserted := []interface{}{}
	order.TableID = orderItemPack.TableID
	orderId, err := OrderItemOrderCreator(r.Context(), order)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	validate := validator.New()
	for _, orderItem := range orderItemPack.OrderItems {
		orderItem.OrderID = orderId
		if err := validate.Struct(orderItem); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Invalid Order Item")
			return
		}

		orderItem.CreatedOn = time.Now()
		orderItem.UpdatedOn = time.Now()
		orderItem.UnitPrice = helpers.ToFixed(orderItem.UnitPrice, 2)

		orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)
	}

	result, err := orderItemCollection.InsertMany(r.Context(), orderItemsToBeInserted)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result.InsertedIDs)
}

func UpdateOrderItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderItemID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid Order Item ID")
		return
	}

	var orderItem models.OrderItem
	if err = json.NewDecoder(r.Body).Decode(&orderItem); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid Order Item")
		return
	}

	validate := validator.New()
	if err := validate.Struct(orderItem); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid Order Item : " + err.Error())
		return
	}

	orderID := orderItem.OrderID
	query := bson.M{
		"_id": orderID,
	}
	count, err := orderCollection.CountDocuments(r.Context(), query)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}
	if count == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Order not found")
		return
	}

	orderItem.UpdatedOn = time.Now()
	orderItem.UnitPrice = helpers.ToFixed(orderItem.UnitPrice, 2)

	query = bson.M{
		"_id": orderItemID,
	}
	update := bson.M{
		"$set": orderItem,
	}

	result, err := orderItemCollection.UpdateOne(r.Context(), query, update)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	if result.MatchedCount == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Order Item not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Order Item updated successfully")
}

func ItemsByOrder(ctx context.Context, orderID primitive.ObjectID) (OrderItems []primitive.M, err error) {
	match := bson.M{
		"$match": bson.M{
			"orderID": orderID,
		},
	}
	lookupFood := bson.M{
		"$lookup": bson.M{
			"from":         "food",
			"localField":   "foodID",
			"foreignField": "_id",
			"as":           "food",
		},
	}
	unwindFood := bson.M{
		"path":                       "$food",
		"preserveNullAndEmptyArrays": true,
	}

	lookupOrder := bson.M{
		"$lookup": bson.M{
			"from":         "order",
			"localField":   "orderID",
			"foreignField": "_id",
			"as":           "order",
		},
	}

	unwindOrder := bson.M{
		"path":                       "$order",
		"preserveNullAndEmptyArrays": true,
	}

	lookupTable := bson.M{
		"$lookup": bson.M{
			"from":         "table",
			"localField":   "order.tableID",
			"foreignField": "_id",
			"as":           "table",
		},
	}

	unwindTable := bson.M{
		"path":                       "$table",
		"preserveNullAndEmptyArrays": true,
	}

	project := bson.M{
		"$project": bson.M{
			"amount":      "food.price",
			"foodID":      1,
			"foodName":    "$food.name",
			"foodImage":   "$food.image",
			"tableNumber": "$table.number",
			"tableID":     "$order.tableID",
			"price":       "$food.price",
			"orderID":     "$order._id",
			"quantity":    1,
		},
	}

	group := bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"orderID":     "$orderID",
				"tableID":     "$tableID",
				"tableNumber": "$tableNumber",
			},
			"paymentDue": bson.M{
				"$sum": bson.M{
					"$multiply": []interface{}{"$price", "$quantity"},
				},
			},
			"totalCount": bson.M{
				"$sum": 1,
			},
			"orderItems": bson.M{
				"$push": "$$ROOT",
			},
		},
	}

	project2 := bson.M{
		"$project": bson.M{
			"_id":         1,
			"paymentDue":  1,
			"orderItems":  1,
			"tableNumber": "$_id.tableNumber",
			"totalCount":  1,
		},
	}

	pipeline := []bson.M{
		match,
		lookupFood,
		unwindFood,
		lookupOrder,
		unwindOrder,
		lookupTable,
		unwindTable,
		project,
		group,
		project2,
	}
	cursor, err := orderItemCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &OrderItems); err != nil {
		return nil, err
	}

	return OrderItems, nil
}
