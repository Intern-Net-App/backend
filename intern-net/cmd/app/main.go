package main

import (
	"context"
	"fmt"
	"intern-net/internal/app/handlers"
	"intern-net/internal/app/repositories"
	"intern-net/internal/app/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Setup MongoDB
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetServerAPIOptions(serverAPI)

	// Create new client and connect to server
	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	userCollection := client.Database(os.Getenv("MONGODB_DATABASE")).Collection("users")
	jobCollection := client.Database(os.Getenv("MONGODB_DATABASE")).Collection("linkedin_jobs")
	userRepository := repositories.NewUserRepository(userCollection)
	jobRepository := repositories.NewJobRepository(jobCollection)

	// Create router using Gorilla Mux
	r := mux.NewRouter()

	// Register handlers and middleware
	r.Handle("/api/signup", handlers.NewSignupHandler(userRepository))
	r.Handle("/api/login", handlers.NewLoginHandler(userRepository))

	// Authentication Middleware to protected routes
	r.Handle("/welcome", services.Authenticate(http.HandlerFunc(handlers.Welcome)))

	// Display Job postings handlers
	r.Handle("/api/jobs", handlers.NewJobPostingsHandler(jobRepository))

	// CORS Handler
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	// Attach CORS Middleware to router
	handler := c.Handler(r)

	// Create HTTP Server
	server := http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
	// Start the server
	go func() {
		log.Println("server listening on :8080")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	// Shut server down on interrupt signal
	stop := make(chan os.Signal, 1) // Change to os.Signal
	go func() {
		sig := <-stop
		fmt.Printf("Received signal: %v\n", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown error: %v", err)
		}
	}()

	// Wait for interrupt
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
	stop <- os.Interrupt // Send the interrupt signal to the correct channel
}
