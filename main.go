package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/toundaherve/sad-api/postgres"
	"github.com/toundaherve/sad-api/user"
)

var Address = ":8001"

func main() {
	validate := validator.New()
	postgresDB := postgres.NewPostgresDB()
	userHandler := user.NewUserHandler(validate, postgresDB)

	router := mux.NewRouter()
	router.Methods("POST").Path("/api/onboarding/begin_verification")
	router.Methods("POST").Path("/api/users").HandlerFunc(userHandler.CreateUser)
	router.Methods("GET").Path("/api/users/email_available").HandlerFunc(userHandler.CheckEmailAvailable)

	srv := http.Server{
		Addr:         Address,
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
