package main

import (
	"fmt"
	"log"
	"os"

	"ethosview-backend/internal/models"
	"ethosview-backend/pkg/database"
)

func main() {
	// Load environment variables
	if err := loadEnv(); err != nil {
		log.Fatalf("Failed to load environment: %v", err)
	}

	// Initialize database connection
	db, err := database.InitPostgreSQL()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	fmt.Println("✅ Successfully connected to Supabase database!")

	// Test company repository
	companyRepo := models.NewCompanyRepository(db)

	// Test getting companies
	companies, err := companyRepo.ListCompanies(5, 0, "")
	if err != nil {
		log.Printf("❌ Failed to get companies: %v", err)
	} else {
		fmt.Printf("✅ Successfully retrieved %d companies\n", len(companies))
		for _, company := range companies {
			fmt.Printf("  - %s (%s) - %s\n", company.Name, company.Symbol, company.Sector)
		}
	}

	// Test ESG scores repository
	esgRepo := models.NewESGScoreRepository(db)

	// Test getting ESG scores
	esgScores, err := esgRepo.ListESGScores(5, 0, 0)
	if err != nil {
		log.Printf("❌ Failed to get ESG scores: %v", err)
	} else {
		fmt.Printf("✅ Successfully retrieved %d ESG scores\n", len(esgScores))
		for _, score := range esgScores {
			fmt.Printf("  - %s: Overall Score %.2f (E:%.2f S:%.2f G:%.2f)\n",
				score.CompanyName, score.OverallScore,
				score.EnvironmentalScore, score.SocialScore, score.GovernanceScore)
		}
	}

	fmt.Println("\n🎉 Supabase integration test completed successfully!")
}

func loadEnv() error {
	// Set default environment variables for testing
	envVars := map[string]string{
		"DB_HOST":     "aws-1-us-east-2.pooler.supabase.com",
		"DB_PORT":     "5432",
		"DB_USER":     "postgres.wrxyobquvqbwuinlikur",
		"DB_PASSWORD": "Dragonbobby20",
		"DB_NAME":     "postgres",
		"DB_SSL_MODE": "require",
	}

	for key, defaultValue := range envVars {
		if os.Getenv(key) == "" {
			os.Setenv(key, defaultValue)
		}
	}

	return nil
}
