package routes

import (
	"github.com/MayankSaxena03/Restaurant-Management-System/controllers"
	"github.com/MayankSaxena03/Restaurant-Management-System/middleware"
	"github.com/gorilla/mux"
)

func OrderItemRoutes(router *mux.Router) {
	OrderItemGroup := router.PathPrefix("/orderItems").Subrouter()
	OrderItemGroup.Use(middleware.Authentication)
	OrderItemGroup.HandleFunc("", controllers.GetOrderItems).Methods("GET")
	OrderItemGroup.HandleFunc("/{id}", controllers.GetOrderItem).Methods("GET")
	OrderItemGroup.HandleFunc("/order/{orderId}", controllers.GetOrderItemsByOrder).Methods("GET")
	OrderItemGroup.HandleFunc("", controllers.CreateOrderItem).Methods("POST")
	OrderItemGroup.HandleFunc("/{id}", controllers.UpdateOrderItem).Methods("PATCH")
}
