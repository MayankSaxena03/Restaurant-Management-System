package routes

import (
	"github.com/MayankSaxena03/Restaurant-Management-System/controllers"
	"github.com/MayankSaxena03/Restaurant-Management-System/middleware"
	"github.com/gorilla/mux"
)

func InvoiceRoutes(router *mux.Router) {
	invoiceGroup := router.PathPrefix("/invoices").Subrouter()
	invoiceGroup.Use(middleware.Authentication)
	invoiceGroup.HandleFunc("", controllers.GetInvoices).Methods("GET")
	invoiceGroup.HandleFunc("/{id}", controllers.GetInvoice).Methods("GET")
	invoiceGroup.HandleFunc("", controllers.CreateInvoice).Methods("POST")
	invoiceGroup.HandleFunc("/{id}", controllers.UpdateInvoice).Methods("PATCH")
}
