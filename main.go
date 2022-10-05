package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"go_jwt_mux/controller/authcontroller"
	"go_jwt_mux/controller/product_controller"
	"go_jwt_mux/middlewares"
	"go_jwt_mux/models"
	"log"
	"net/http"
)

func main() {
	models.ConnectionDatabase()

	r := mux.NewRouter()

	r.HandleFunc("/login", authcontroller.Login).Methods(http.MethodPost)
	r.HandleFunc("/register", authcontroller.Register).Methods(http.MethodPost)
	r.HandleFunc("/logout", authcontroller.Logout).Methods(http.MethodGet)

	// grouping endpoint in /api
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/products", product_controller.Index).Methods(http.MethodGet)
	api.Use(middlewares.JWTMiddleware)

	fmt.Println("Server Run And Listening at Port 8010..")
	log.Fatal(http.ListenAndServe(":8010", r))
}
