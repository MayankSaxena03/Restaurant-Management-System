package routes

import (
	"github.com/MayankSaxena03/Restaurant-Management-System/controllers"
	"github.com/MayankSaxena03/Restaurant-Management-System/middleware"
	"github.com/gorilla/mux"
)

func TableRoutes(router *mux.Router) {
	tableGroup := router.PathPrefix("/tables").Subrouter()
	tableGroup.Use(middleware.Authentication)
	tableGroup.HandleFunc("", controllers.GetTables).Methods("GET")
	tableGroup.HandleFunc("/{id}", controllers.GetTable).Methods("GET")
	tableGroup.HandleFunc("", controllers.CreateTable).Methods("POST")
	tableGroup.HandleFunc("/{id}", controllers.UpdateTable).Methods("PATCH")
}
