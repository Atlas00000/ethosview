package server

import (
	"database/sql"
	"net/http"
	"time"

	"ethosview-backend/internal/handlers"
	"ethosview-backend/internal/websocket"
	"ethosview-backend/pkg/auth"
	"ethosview-backend/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Server represents the HTTP server
type Server struct {
	router    *gin.Engine
	db        *sql.DB
	redis     *redis.Client
	wsManager *websocket.Manager
}

// NewServer creates and configures a new server instance
func NewServer(db *sql.DB, redis *redis.Client) *Server {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create router
	router := gin.Default()

	// Create server instance
	srv := &Server{
		router:    router,
		db:        db,
		redis:     redis,
		wsManager: websocket.NewManager(),
	}

	// Setup routes
	srv.setupRoutes()

	return srv
}

// setupRoutes configures all the routes for the server
func (s *Server) setupRoutes() {
	// Initialize performance middleware
	rateLimiter := middleware.NewRateLimiter(s.redis)
	monitoringMiddleware := middleware.MonitoringMiddleware()
	cacheMiddleware := middleware.CacheMiddleware(s.redis, 5*time.Minute)

	// Apply global middleware
	s.router.Use(monitoringMiddleware)

	// Health check endpoint
	s.router.GET("/health", s.healthCheck)
	s.router.GET("/metrics", middleware.HealthCheckMiddleware())

	// Initialize JWT manager and middleware
	jwtManager := auth.NewJWTManager()
	authMiddleware := middleware.AuthMiddleware(jwtManager)

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// Health check for API
		v1.GET("/health", s.healthCheck)

		// Initialize handlers
		authHandler := handlers.NewAuthHandler(s.db)
		companyHandler := handlers.NewCompanyHandler(s.db)
		esgHandler := handlers.NewESGHandler(s.db)
		dashboardHandler := handlers.NewDashboardHandler(s.db)
		financialHandler := handlers.NewFinancialHandler(s.db)
		analyticsHandler := handlers.NewAnalyticsHandler(s.db)
		advancedAnalyticsHandler := handlers.NewAdvancedAnalyticsHandler(s.db)

		// Authentication routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/profile", authMiddleware, authHandler.GetProfile)
			auth.PUT("/profile", authMiddleware, authHandler.UpdateProfile)
		}

		// Company routes (public for now, can be protected later)
		companies := v1.Group("/companies")
		companies.Use(rateLimiter.RateLimitMiddleware(100)) // 100 requests per minute
		companies.Use(cacheMiddleware)
		{
			companies.POST("", companyHandler.CreateCompany)
			companies.GET("", companyHandler.ListCompanies)
			companies.GET("/sectors", companyHandler.GetSectors)
			companies.GET("/symbol/:symbol", companyHandler.GetCompanyBySymbol)
			companies.GET("/:id", companyHandler.GetCompany)
			companies.PUT("/:id", companyHandler.UpdateCompany)
			companies.DELETE("/:id", companyHandler.DeleteCompany)
		}

		// ESG routes (public for now, can be protected later)
		esg := v1.Group("/esg")
		{
			esg.POST("/scores", esgHandler.CreateESGScore)
			esg.GET("/scores", esgHandler.ListESGScores)
			esg.GET("/scores/:id", esgHandler.GetESGScore)
			esg.PUT("/scores/:id", esgHandler.UpdateESGScore)
			esg.DELETE("/scores/:id", esgHandler.DeleteESGScore)
			esg.GET("/companies/:id/latest", esgHandler.GetLatestESGScoreByCompany)
			esg.GET("/companies/:id/scores", esgHandler.GetESGScoresByCompany)
		}

		// Dashboard route
		v1.GET("/dashboard", dashboardHandler.GetDashboard)

		// Financial routes (public for now, can be protected later)
		financial := v1.Group("/financial")
		{
			financial.GET("/companies/:id/prices", financialHandler.GetStockPrices)
			financial.GET("/companies/:id/price/latest", financialHandler.GetLatestStockPrice)
			financial.GET("/companies/:id/indicators", financialHandler.GetFinancialIndicators)
			financial.GET("/companies/:id/summary", financialHandler.GetCompanyFinancialSummary)
			financial.GET("/market", financialHandler.GetMarketData)
			financial.GET("/market/history", financialHandler.GetMarketDataHistory)
		}

		// Analytics routes (public for now, can be protected later)
		analytics := v1.Group("/analytics")
		analytics.Use(rateLimiter.RateLimitMiddleware(50)) // 50 requests per minute for analytics
		analytics.Use(cacheMiddleware)
		{
			analytics.GET("/companies/:id/esg-trends", analyticsHandler.GetESGTrends)
			analytics.GET("/sectors/comparisons", analyticsHandler.GetSectorComparisons)
			analytics.GET("/financial/comparisons", analyticsHandler.GetFinancialComparisons)
			analytics.GET("/top-performers/:metric", analyticsHandler.GetTopPerformers)
			analytics.GET("/correlation/esg-financial", analyticsHandler.GetESGvsFinancialCorrelation)
			analytics.GET("/summary", analyticsHandler.GetAnalyticsSummary)
		}

		// Advanced Analytics routes (public for now, can be protected later)
		advanced := v1.Group("/advanced")
		advanced.Use(rateLimiter.RateLimitMiddleware(30)) // 30 requests per minute for advanced analytics
		advanced.Use(cacheMiddleware)
		{
			advanced.GET("/companies/:id/predict-esg", advancedAnalyticsHandler.PredictESGScore)
			advanced.GET("/portfolio/optimize", advancedAnalyticsHandler.OptimizePortfolio)
			advanced.GET("/companies/:id/risk-assessment", advancedAnalyticsHandler.AssessRisk)
			advanced.GET("/companies/:id/trends/:metric", advancedAnalyticsHandler.AnalyzeTrend)
			advanced.GET("/summary", advancedAnalyticsHandler.GetAdvancedAnalyticsSummary)
		}

		// WebSocket routes
		wsHandler := handlers.NewWebSocketHandler(s.wsManager)
		v1.GET("/ws", wsHandler.HandleWebSocket)
		v1.GET("/ws/status", wsHandler.GetWebSocketStatus)
	}
}

// healthCheck handles health check requests
func (s *Server) healthCheck(c *gin.Context) {
	// Check database connection
	if err := s.db.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": "Database connection failed",
			"error":   err.Error(),
		})
		return
	}

	// Check Redis connection
	ctx := c.Request.Context()
	if err := s.redis.Ping(ctx).Err(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": "Redis connection failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"message": "All services are running",
		"version": "1.0.0",
	})
}

// Run starts the HTTP server
func (s *Server) Run(addr string) error {
	// Start WebSocket manager in a goroutine
	go s.wsManager.Start()

	return s.router.Run(addr)
}
