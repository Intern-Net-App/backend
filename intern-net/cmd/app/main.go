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

	"github.com/joho/godotenv"
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
	userRepository := repositories.NewUserRepository(userCollection)

	// Create HTTP Server
	server := http.Server{
		Addr: ":8080",
	}

	// Register handlers and middleware
	http.Handle("/api/signup", handlers.NewSignupHandler(userRepository))
	http.Handle("/api/login", handlers.NewLoginHandler(userRepository))

	// Authentication Middleware to protected routes
	http.Handle("/welcome", services.Authenticate(http.HandlerFunc(handlers.Welcome)))

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
