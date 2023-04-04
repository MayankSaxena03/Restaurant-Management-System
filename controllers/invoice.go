package controllers

import (
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

type InvoiceViewFormat struct {
	ID             string      `json:"_id,omitempty" bson:"_id,omitempty"`
	OrderID        string      `json:"orderID,omitempty" bson:"orderID,omitempty"`
	OrderDetails   interface{} `json:"orderDetails,omitempty" bson:"orderDetails,omitempty"`
	TableID        string      `json:"tableNumber,omitempty" bson:"tableNumber,omitempty"`
	PaymentMethod  string      `json:"paymentMethod,omitempty" bson:"paymentMethod,omitempty"`
	PaymentStatus  string      `json:"paymentStatus,omitempty" bson:"paymentStatus,omitempty"`
	PaymentDue     interface{} `json:"paymentDue,omitempty" bson:"paymentDue,omitempty"`
	PaymentDueDate time.Time   `json:"paymentDueDate,omitempty" bson:"paymentDueDate,omitempty"`
}

var invoiceCollection = database.OpenCollection(database.Client, "invoice")

func GetInvoices(w http.ResponseWriter, r *http.Request) {
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

	var invoices []models.Invoice
	cursor, err := invoiceCollection.Find(r.Context(), bson.M{}, options.Find().SetLimit(int64(limit)).SetSkip(int64(skip)))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}
	if err = cursor.All(r.Context(), &invoices); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoices)
}

func GetInvoice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if !primitive.IsValidObjectID(id) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid ID")
		return
	}

	invoiceID, _ := primitive.ObjectIDFromHex(id)
	var invoice models.Invoice
	err := invoiceCollection.FindOne(r.Context(), bson.M{"_id": invoiceID}).Decode(&invoice)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode("Invoice not found")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	var invoiceView InvoiceViewFormat
	allOrderItems, err := ItemsByOrder(r.Context(), invoice.OrderID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}
	invoiceView.OrderID = invoice.OrderID.Hex()
	invoiceView.PaymentDueDate = invoice.PaymentDueDate
	invoiceView.PaymentMethod = invoice.PaymentMethod
	invoiceView.ID = invoice.ID.Hex()
	invoiceView.PaymentStatus = invoice.PaymentStatus
	invoiceView.PaymentDue = allOrderItems[0]["paymentDue"]
	invoiceView.TableID = allOrderItems[0]["tableID"].(string)
	invoiceView.OrderDetails = allOrderItems[0]["orderItems"]

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invoiceView)
}

func CreateInvoice(w http.ResponseWriter, r *http.Request) {
	var invoice models.Invoice
	err := json.NewDecoder(r.Body).Decode(&invoice)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid request body")
		return
	}

	validate := validator.New()
	err = validate.Struct(invoice)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid request body : " + err.Error())
		return
	}

	count, err := orderCollection.CountDocuments(r.Context(), bson.M{"_id": invoice.OrderID})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}
	if count == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Order not found")
		return
	}

	invoice.CreatedOn = time.Now()
	invoice.UpdatedOn = time.Now()
	invoice.PaymentDueDate = invoice.PaymentDueDate.AddDate(0, 0, 1)
	result, err := invoiceCollection.InsertOne(r.Context(), invoice)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	invoice.ID = result.InsertedID.(primitive.ObjectID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(invoice)
}

func UpdateInvoice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if !primitive.IsValidObjectID(id) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid ID")
		return
	}

	var updatedInvoice models.Invoice
	err := json.NewDecoder(r.Body).Decode(&updatedInvoice)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid request body")
		return
	}

	if !primitive.IsValidObjectID(updatedInvoice.OrderID.Hex()) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Invalid Order ID")
		return
	}

	count, err := orderCollection.CountDocuments(r.Context(), bson.M{"_id": updatedInvoice.OrderID})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}
	if count == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Order not found")
		return
	}

	if updatedInvoice.PaymentMethod == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Payment method is required")
		return
	}

	if updatedInvoice.PaymentStatus == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Payment status is required")
		return
	}

	if updatedInvoice.PaymentDueDate.IsZero() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Payment due date is required")
		return
	}

	invoiceID, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{
		"_id": invoiceID,
	}
	update := bson.M{
		"$set": bson.M{
			"paymentMethod":  updatedInvoice.PaymentMethod,
			"paymentStatus":  updatedInvoice.PaymentStatus,
			"paymentDueDate": updatedInvoice.PaymentDueDate,
			"orderID":        updatedInvoice.OrderID,
			"updatedOn":      time.Now(),
		},
	}

	_, err = invoiceCollection.UpdateOne(r.Context(), query, update)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Something went wrong")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Invoice updated successfully")
}
