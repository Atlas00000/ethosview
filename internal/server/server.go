package server

import (
	"database/sql"
	"net/http"
	"time"

	"ethosview-backend/internal/handlers"
	"ethosview-backend/internal/websocket"
	"ethosview-backend/pkg/auth"
	"ethosview-backend/pkg/cache"
	"ethosview-backend/pkg/dashboard"
	"ethosview-backend/pkg/health"
	"ethosview-backend/pkg/metrics"
	"ethosview-backend/pkg/middleware"
	"ethosview-backend/pkg/monitoring"
	"ethosview-backend/pkg/security"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Server represents the HTTP server
type Server struct {
	router             *gin.Engine
	db                 *sql.DB
	redis              *redis.Client
	wsManager          *websocket.Manager
	cacheWarmer        *cache.CacheWarmer
	advancedCache      *cache.AdvancedCache
	metricsCollector   *metrics.MetricsCollector
	healthChecker      *health.HealthChecker
	securityMiddleware *security.SecurityMiddleware
	businessDashboard  *dashboard.BusinessDashboard
	alertManager       *monitoring.AlertManager
}

// NewServer creates and configures a new server instance
func NewServer(db *sql.DB, redis *redis.Client) *Server {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create router
	router := gin.Default()

	// Create server instance
	srv := &Server{
		router:             router,
		db:                 db,
		redis:              redis,
		wsManager:          websocket.NewManager(),
		cacheWarmer:        cache.NewCacheWarmer(redis, db),
		advancedCache:      cache.NewAdvancedCache(redis, "ethosview"),
		metricsCollector:   metrics.NewMetricsCollector(redis, db),
		healthChecker:      health.NewHealthChecker(db, redis),
		securityMiddleware: security.NewSecurityMiddleware(),
		businessDashboard:  dashboard.NewBusinessDashboard(db, redis),
		alertManager:       monitoring.NewAlertManager(db, redis),
	}

	// Setup routes
	srv.setupRoutes()

	// Start background services
	srv.startBackgroundServices()

	return srv
}

// setupRoutes configures all the routes for the server
func (s *Server) setupRoutes() {
	// Initialize performance middleware
	rateLimiter := middleware.NewRateLimiter(s.redis)
	monitoringMiddleware := middleware.MonitoringMiddleware()
	cacheMiddleware := middleware.CacheMiddleware(s.redis, 5*time.Minute)
	compressionMiddleware := middleware.CompressionMiddleware()
	requestIDMiddleware := middleware.RequestIDMiddleware()

	// Apply global middleware
	s.router.Use(requestIDMiddleware)
	s.router.Use(compressionMiddleware)
	s.router.Use(monitoringMiddleware)
	s.router.Use(s.securityMiddleware.SecurityHeaders())
	s.router.Use(s.securityMiddleware.CORS())
	s.router.Use(s.securityMiddleware.InputSanitization())
	s.router.Use(s.securityMiddleware.SQLInjectionProtection())
	s.router.Use(s.securityMiddleware.XSSProtection())
	s.router.Use(s.securityMiddleware.RequestSizeLimit(10 * 1024 * 1024)) // 10MB limit

	// Health check endpoints
	s.router.GET("/health", s.healthChecker.HealthCheckHandler())
	s.router.GET("/health/detailed", s.healthChecker.DetailedHealthCheckHandler())
	s.router.GET("/health/ready", s.healthChecker.ReadinessCheckHandler())
	s.router.GET("/health/live", s.healthChecker.LivenessCheckHandler())
	s.router.GET("/metrics", s.metricsHandler)
	s.router.GET("/alerts", s.alertsHandler)
	s.router.GET("/dashboard/business", s.businessDashboardHandler)

	// Initialize JWT manager and middleware
	jwtManager := auth.NewJWTManager()
	authMiddleware := middleware.AuthMiddleware(jwtManager)

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// Health check for API
		v1.GET("/health", s.healthChecker.HealthCheckHandler())

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

// startBackgroundServices starts background services
func (s *Server) startBackgroundServices() {
	// Start cache warming (every 30 minutes)
	s.cacheWarmer.StartCacheWarming(30 * time.Minute)

	// Start metrics collection (every 5 minutes)
	s.metricsCollector.StartMetricsCollection(5 * time.Minute)
	
	// Start performance monitoring and alerting (every 1 minute)
	s.alertManager.StartMonitoring(1 * time.Minute)
}

// metricsHandler handles metrics requests
func (s *Server) metricsHandler(c *gin.Context) {
	metrics := s.metricsCollector.GetMetrics()
	c.JSON(http.StatusOK, metrics)
}

// alertsHandler handles alerts requests
func (s *Server) alertsHandler(c *gin.Context) {
	// Check if requesting all alerts or just active ones
	showAll := c.Query("all") == "true"
	
	var alerts []monitoring.Alert
	if showAll {
		alerts = s.alertManager.GetAllAlerts()
	} else {
		alerts = s.alertManager.GetActiveAlerts()
	}
	
	c.JSON(http.StatusOK, gin.H{
		"alerts": alerts,
		"count":  len(alerts),
		"active_only": !showAll,
	})
}

// businessDashboardHandler handles business dashboard requests
func (s *Server) businessDashboardHandler(c *gin.Context) {
	dashboard, err := s.businessDashboard.GetDashboardData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve dashboard data",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, dashboard)
}

// Run starts the HTTP server
func (s *Server) Run(addr string) error {
	// Start WebSocket manager in a goroutine
	go s.wsManager.Start()

	return s.router.Run(addr)
}
