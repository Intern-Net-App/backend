// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/gorilla/mux"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"

// 	"intern-net/internal/app/auth"
// 	"intern-net/internal/app/handlers"
// 	"intern-net/internal/app/repositories"
// )

// func main() {
// 	// Set up MongoDB connection
// 	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
// 	client, err := mongo.Connect(context.TODO(), clientOptions)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer client.Disconnect(context.TODO())

// 	// Access the MongoDB collection
// 	collection := client.Database(os.Getenv("MONGODB_NAME")).Collection("users")

// 	// Initialize the user repository with the collection
// 	userRepo := repositories.NewUserRepository(collection)

// 	// Create a new router
// 	router := mux.NewRouter()

// 	// Define routes

// 	// Route for user login or signup
// 	router.HandleFunc("/api/login", handlers.LoginHandler(userRepo)).Methods("POST")

// 	// Protected route - example
// 	router.Handle("/api/protected", auth.Middleware(http.HandlerFunc(handlers.ProtectedHandler))).Methods("GET")

// 	// Start the HTTP server
// 	port := 8080
// 	fmt.Printf("Server is running on :%d\n", port)
// 	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
// }

package main

import (
	"intern-net/internal/app/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/welcome", handlers.Welcome)
	http.HandleFunc("/refresh", handlers.Refresh)
	http.HandleFunc("/logout", handlers.Logout)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
