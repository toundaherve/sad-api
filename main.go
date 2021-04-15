package main

import (
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/toundaherve/sad-api/csvfile"
	"github.com/toundaherve/sad-api/logger"
	"github.com/toundaherve/sad-api/onboarding"
	"github.com/toundaherve/sad-api/user"
)

var Address = ":8001"
var csvFilename = "users.csv"

func main() {
	validate := validator.New()
	logger := logger.NewLogger()
	// postgresDB := postgres.NewPostgresDB()
	csvFile := csvfile.New(csvFilename)
	// if err := csvFile.InitFile(); err != nil {
	// 	logger.WithField("filename", csvFilename).Fatalln("Failed to create csvFile Storage")
	// }
	onboardingHandler := onboarding.New(nil, logger)
	userHandler := user.NewUserHandler(validate, csvFile, logger)

	router := mux.NewRouter()
	router.Methods("GET").Path("/api/onboarding/begin_verification").HandlerFunc(onboardingHandler.BeginVerification)
	router.Methods("POST").Path("/api/onboarding/verify_code").HandlerFunc(onboardingHandler.VerifyCode)
	router.Methods("POST").Path("/api/users").HandlerFunc(userHandler.CreateUser)
	router.Methods("GET").Path("/api/users/email_available").HandlerFunc(userHandler.CheckEmailAvailable)
	router.Methods("POST").Path("/api/users/login").HandlerFunc(userHandler.Login)

	handler := cors.Default().Handler(router)

	srv := http.Server{
		Handler:      handler,
		Addr:         Address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.WithField("err", err.Error()).Fatalln("Failed to start server")
	}
}

// func CORS(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Access-Control-Allow-Headers", "*")
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Allow-Methods", "*")

// 		if r.Method == "OPTIONS" {
// 			w.WriteHeader(http.StatusOK)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }
