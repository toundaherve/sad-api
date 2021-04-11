package main

import (
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/toundaherve/sad-api/logger"
	"github.com/toundaherve/sad-api/postgres"
	"github.com/toundaherve/sad-api/user"
)

var Address = ":8001"

func main() {
	validate := validator.New()
	postgresDB := postgres.NewPostgresDB()
	logger := logger.NewLogger()
	userHandler := user.NewUserHandler(validate, postgresDB, logger)

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

	if err := srv.ListenAndServe(); err != nil {
		logger.WithField("err", err.Error()).Fatalln("Failed to start server")
	}
}
