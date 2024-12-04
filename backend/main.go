package main

import (
	"log"
	"net/http"
	"time"

	"github.com/UpsDev42069/BM_Search_Engine/backend/config"
	"github.com/UpsDev42069/BM_Search_Engine/backend/db"
	"github.com/UpsDev42069/BM_Search_Engine/backend/handlers"
	"github.com/UpsDev42069/BM_Search_Engine/backend/metrics"
	"github.com/UpsDev42069/BM_Search_Engine/backend/security"
	"github.com/rs/cors"

	_ "github.com/UpsDev42069/BM_Search_Engine/backend/docs"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title BM Search Engine API
// @version 2.0
// @description This is a sample server for a BM Search Engine.
// @host localhost:8080
// @BasePath /
func main() {
	// Load environment variables
	config.LoadEnv()

	db.InitializeDBenv()

	security.InitializeStore(config.SessionSecret)

	// Init metrics
	metrics.Init()

	go metrics.CollectSystemMetrics(10 * time.Second)

	// Connecting to the database
	database, err := db.ConnectDB(false)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	db.RunMigrations(database)

	r := mux.NewRouter()

	// Middleware for metrics
	r.Use(metrics.Middleware)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{config.FrontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(r)

	// Existing routes
	r.HandleFunc("/", handlers.RootGet).Methods("GET")
	r.HandleFunc("/api/search", handlers.SearchHandler(database)).Methods("GET")
	r.HandleFunc("/api/register", handlers.RegisterHandler(database)).Methods("POST")
	r.HandleFunc("/api/login", handlers.LoginHandler(database)).Methods("POST")
	r.HandleFunc("/api/weather", handlers.WeatherHandler).Methods("GET")
	r.HandleFunc("/api/logout", handlers.LogoutHandler).Methods("GET")
	r.HandleFunc("/api/reset-password", handlers.ResetPasswordHandler(database)).Methods("PUT")
	r.HandleFunc("/api/check-login", handlers.CheckLoginHandler).Methods("GET")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
