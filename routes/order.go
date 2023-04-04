package routes

import (
	"github.com/MayankSaxena03/Restaurant-Management-System/controllers"
	"github.com/MayankSaxena03/Restaurant-Management-System/middleware"
	"github.com/gorilla/mux"
)

func OrderRoutes(router *mux.Router) {
	orderGroup := router.PathPrefix("/orders").Subrouter()
	orderGroup.Use(middleware.Authentication)
	orderGroup.HandleFunc("", controllers.GetOrders).Methods("GET")
	orderGroup.HandleFunc("/{id}", controllers.GetOrder).Methods("GET")
	orderGroup.HandleFunc("", controllers.CreateOrder).Methods("POST")
	orderGroup.HandleFunc("/{id}", controllers.UpdateOrder).Methods("PATCH")
}
