package routes

import (
	"github.com/MayankSaxena03/Restaurant-Management-System/controllers"
	"github.com/MayankSaxena03/Restaurant-Management-System/middleware"
	"github.com/gorilla/mux"
)

func MenuRoutes(router *mux.Router) {
	menuGroup := router.PathPrefix("/menus").Subrouter()
	menuGroup.Use(middleware.Authentication)
	menuGroup.HandleFunc("", controllers.GetMenus).Methods("GET")
	menuGroup.HandleFunc("/{id}", controllers.GetMenu).Methods("GET")
	menuGroup.HandleFunc("", controllers.CreateMenu).Methods("POST")
	menuGroup.HandleFunc("/{id}", controllers.UpdateMenu).Methods("PATCH")
}
