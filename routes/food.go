package routes

import (
	"github.com/MayankSaxena03/Restaurant-Management-System/controllers"
	"github.com/MayankSaxena03/Restaurant-Management-System/middleware"
	"github.com/gorilla/mux"
)

func FoodRoutes(router *mux.Router) {
	foodGroup := router.PathPrefix("/foods").Subrouter()
	foodGroup.Use(middleware.Authentication)
	foodGroup.HandleFunc("", controllers.GetFoods).Methods("GET")
	foodGroup.HandleFunc("/{id}", controllers.GetFood).Methods("GET")
	foodGroup.HandleFunc("", controllers.CreateFood).Methods("POST")
	foodGroup.HandleFunc("/{id}", controllers.UpdateFood).Methods("PATCH")
}
