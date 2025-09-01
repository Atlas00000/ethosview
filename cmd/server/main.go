package main

import (
	"log"
	"os"

	"ethosview-backend/internal/server"
	"ethosview-backend/pkg/database"
)

func main() {
	// Initialize database connections
	db, err := database.InitPostgreSQL()
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}
	defer db.Close()

	redisClient, err := database.InitRedis()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer redisClient.Close()

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize and start server
	srv := server.NewServer(db, redisClient)
	log.Printf("Starting server on port %s", port)
	if err := srv.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
