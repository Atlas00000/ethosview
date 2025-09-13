package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables
	loadEnv()

	// Create connection string
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"),
	)

	fmt.Printf("🔗 Testing connection to: %s\n", os.Getenv("DB_HOST"))

	// Open database connection
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("❌ Failed to open database: %v", err)
	}
	defer db.Close()

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Failed to ping database: %v", err)
	}

	fmt.Println("✅ Successfully connected to Supabase database!")

	// Test a simple query
	var version string
	err = db.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		log.Printf("❌ Failed to get version: %v", err)
	} else {
		fmt.Printf("📊 Database version: %s\n", version)
	}

	// Check if our tables exist
	tables := []string{"users", "companies", "esg_scores"}
	for _, table := range tables {
		var exists bool
		query := "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = $1)"
		err = db.QueryRow(query, table).Scan(&exists)
		if err != nil {
			log.Printf("❌ Failed to check table %s: %v", table, err)
		} else if exists {
			fmt.Printf("✅ Table '%s' exists\n", table)
		} else {
			fmt.Printf("⚠️  Table '%s' does not exist (need to run migration)\n", table)
		}
	}

	fmt.Println("\n🎉 Connection test completed!")
}

func loadEnv() {
	// Set default environment variables for testing
	envVars := map[string]string{
		"DB_HOST":     "db.wryyobquvqbwuinkikur.supabase.co",
		"DB_PORT":     "5432",
		"DB_USER":     "postgres",
		"DB_PASSWORD": "Dragonbobby20",
		"DB_NAME":     "postgres",
		"DB_SSL_MODE": "require",
	}

	for key, defaultValue := range envVars {
		if os.Getenv(key) == "" {
			os.Setenv(key, defaultValue)
		}
	}
}
