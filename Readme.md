# **Restaurant Management Backend**

This is a backend implementation for a restaurant management system, written in Go. It provides RESTful API endpoints for managing various aspects of a restaurant, including tables, menus, orders, and invoices.

Project Structure


```
├── controllers
│   ├── food.go
│   ├── invoice.go
│   ├── menu.go
│   ├── order.go
│   ├── orderItem.go
│   ├── table.go
│   └── user.go
├── database
│   └── connection.go
├── helpers
│   ├── costEstimate.go
│   ├── password.go
│   └── token.go
├── middleware
│   └── auth.go
├── models
│   ├── food.go
│   ├── invoice.go
│   ├── menu.go
│   ├── order.go
│   ├── orderItem.go
│   ├── table.go
│   └── user.go
├── routes
│   ├── food_routes.go
│   ├── invoice_routes.go
│   ├── menu_routes.go
│   ├── order_routes.go
│   ├── order_item_routes.go
│   ├── table_routes.go
│   ├── user_auth_routes.go
│   └── user_routes.go
├── .env
├── go.mod
├── go.sum
└── main.go
```

# Description of Directories and Files

    controllers: This directory contains the controller functions for handling HTTP requests for various resources such as foods, invoices, menus, orders, order items, tables, and users.
    database: This directory contains a single file connection.go that sets up a connection to the MongoDB database.
    helpers: This directory contains helper functions for handling tasks such as estimating the cost of an order, hashing and verifying passwords, and generating authentication tokens.
    middleware: This directory contains middleware functions for implementing authentication and authorization for protected routes.
    models: This directory contains the Go structs that define the schema for the various resources in the MongoDB database.
    routes: This directory contains the router functions that define the HTTP routes for the various resources and link them to the appropriate controller functions.
    .env: This file contains the environment variables needed to connect to the MongoDB database and to generate and verify authentication tokens.
    go.mod and go.sum: These files are used by Go modules to manage the project's dependencies.
    main.go: This file contains the main function that sets up the HTTP server and the router.

# API Endpoints

The API endpoints are organized by resource type and are protected by authentication using JWT tokens.

User Authentication

    POST /login: User login endpoint.
    POST /signup: User registration endpoint.

Food

    GET /foods: Get all the food items.
    GET /foods/{id}: Get a specific food item by its ID.
    POST /foods: Create a new food item.
    PATCH /foods/{id}: Update an existing food item.

Invoice

    GET /invoices: Get all the invoices.
    GET /invoices/{id}: Get a specific invoice by its ID.
    POST /invoices: Create a new invoice.
    PATCH /invoices/{id}: Update an existing invoice.

Menu

    GET /menus: Get all the menus.
    GET /menus/{id}: Get a specific menu by its ID.
    POST /menus: Create a new menu.
    PATCH /menus/{id}: Update an existing menu.

Order

    GET /orders: Get all the orders.
    GET /orders/{id}: Get a specific order by its ID.
    POST /orders: Create a new order.
    PATCH /orders/{id}: Update an existing order.

Order Item

    GET /orderItems: Get all the order items.
    GET /orderItems/{id}: Get a specific order item by its ID.
    GET /orderItems/order/{orderId}: Get all the order items associated with a specific order.
    POST /orderItems: Create a new order item.
    PATCH /orderItems/{id}: Update an existing order item.

Table

    GET /tables: Get all the tables.
    GET /tables/{id}: Get a specific table by its ID.
    POST /tables: Create a new table.
    PATCH /tables/{id}: Update an existing table.

User

    GET /users: Get all the users.
    GET /users/{id}: Get a specific user by its ID.

Middleware

    Authentication middleware: Used to authenticate a user before accessing any protected endpoints.

# Project Structure

Here is the structure of the project's file system:

    controllers/: Contains the application logic for each endpoint.
    database/: Contains the file for database connection.
    helpers/: Contains the helper functions used in the application.
    middleware/: Contains the middleware functions used in the application.
    models/: Contains the data models used in the application.
    routes/: Contains the route handlers for each endpoint.
    main.go: The entry point of the application.
    .env: Contains the environment variables used in the application.
    go.mod: Contains the dependencies used in the application.
    go.sum: Contains the checksums for the dependencies used in the application.

Dependencies

Here is the list of dependencies used in this project:

    github.com/gorilla/mux: A powerful HTTP router and URL matcher for building Go web servers.
    go.mongodb.org/mongo-driver/mongo: The official MongoDB driver for Go.
    github.com/joho/godotenv: A Go library for loading environment variables from a .env file.

# Usage
To use this restaurant management backend project, follow these steps:

Clone the repository to your local machine:

    git clone https://github.com/MayankSaxena03/Restaurant-Management-System.git

Install the required dependencies:


    cd Restaurant-Management-System
    go mod tidy

Create a .env file in the root directory and add your configuration details in the following format:

    MONGO_URI=<your-mongodb-uri>
    SECRETKEY=<your-jwt-secret>

Start the server:

    go run main.go

Use a tool like Postman to interact with the REST API endpoints. The base URL for the API is http://localhost:8080/.

Create a new user by sending a POST request to /signup endpoint with the required user details. Once the user is created, use the /login endpoint to get an access token.

Use the obtained access token in the Authorization header (Bearer Token) for all protected routes that require authentication.

You can use the provided API endpoints to perform CRUD operations on tables, food items, orders, and invoices. Use the following endpoints as per your needs:

`/tables`: GET all tables, GET a specific table by ID, CREATE a new table, UPDATE an existing table

`/foods`: GET all food items, GET a specific food item by ID, CREATE a new food item, UPDATE an existing food item

`/menus`: GET all menus, GET a specific menu by ID, CREATE a new menu, UPDATE an existing menu

`/orders`: GET all orders, GET a specific order by ID, CREATE a new order, UPDATE an existing order

`/orderItems`: GET all order items, GET a specific order item by ID, GET all order items for a specific order, CREATE a new order item, UPDATE an existing order item

`/invoices`: GET all invoices, GET a specific invoice by ID, CREATE a new invoice, UPDATE an existing invoice
