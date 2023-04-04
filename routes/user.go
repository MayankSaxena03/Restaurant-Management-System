package routes

import (
	"github.com/MayankSaxena03/Restaurant-Management-System/controllers"
	"github.com/MayankSaxena03/Restaurant-Management-System/middleware"
	"github.com/gorilla/mux"
)

func UserAuthRoutes(router *mux.Router) {
	router.HandleFunc("/login", controllers.Login).Methods("POST")
	router.HandleFunc("/signup", controllers.SignUp).Methods("POST")
}

func UserRoutes(router *mux.Router) {
	userGroup := router.PathPrefix("/users").Subrouter()
	userGroup.Use(middleware.Authentication)
	userGroup.HandleFunc("", controllers.GetUsers).Methods("GET")
	userGroup.HandleFunc("/{id}", controllers.GetUser).Methods("GET")
}
