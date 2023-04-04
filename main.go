package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MayankSaxena03/Restaurant-Management-System/routes"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := mux.NewRouter()

	routes.UserAuthRoutes(router)
	// router.MiddlewareFunc(middleware.Authentication)
	routes.FoodRoutes(router)
	routes.InvoiceRoutes(router)
	routes.MenuRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.TableRoutes(router)
	routes.UserRoutes(router)

	fmt.Println("Server is running on port " + port)
	http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, router)) // LoggingHandler is used to log all the requests
}
