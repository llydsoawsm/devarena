package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"devcode/backend/internal/config"
	"devcode/backend/internal/database"
	"devcode/backend/internal/handlers"
	"devcode/backend/internal/middleware"
	"devcode/backend/internal/services"
	"devcode/backend/internal/ws"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	log.Printf("Starting DevCode server in %s mode on port %s", cfg.AppEnv, cfg.Port)

	// Database setup
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	// Seed algorithmic problems
	if err := services.SeedProblems(db); err != nil {
		log.Printf("Seeding warning: %v", err)
	}

	// Real-time WebSocket Hub
	hub := ws.InitHub()

	// Judge Sandbox Runner
	judge := services.NewJudgeService(cfg.JudgeTimeoutSeconds)

	// Matchmaking Engine
	engine := services.NewMatchmakingEngine(db, cfg.RedisAddr, cfg.RedisPassword, hub)

	// HTTP Handlers
	authHandler := handlers.NewAuthHandler(db, cfg)
	problemHandler := handlers.NewProblemHandler(db)
	matchmakingHandler := handlers.NewMatchmakingHandler(db, engine)
	matchHandler := handlers.NewMatchHandler(db, engine)
	submissionHandler := handlers.NewSubmissionHandler(db, judge, engine, hub)

	// Gin Router
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// CORS Configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	r.Use(cors.New(corsConfig))

	// Health Check
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
		})
	})

	// Public Auth Endpoints
	api := r.Group("/api")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh", authHandler.Refresh)
			authGroup.GET("/me", middleware.AuthMiddleware(cfg.JWTSecret), authHandler.Me)
		}

		// Problems
		problemGroup := api.Group("/problems")
		{
			problemGroup.GET("", problemHandler.List)
			problemGroup.GET("/:slug", problemHandler.GetBySlug)
			problemGroup.GET("/id/:id", problemHandler.GetByID)
		}

		// Matchmaking (Protected)
		matchmakingGroup := api.Group("/matchmaking", middleware.AuthMiddleware(cfg.JWTSecret))
		{
			matchmakingGroup.POST("/enqueue", matchmakingHandler.Enqueue)
			matchmakingGroup.POST("/dequeue", matchmakingHandler.Dequeue)
			matchmakingGroup.GET("/status", matchmakingHandler.Status)
		}

		// Matches
		matchGroup := api.Group("/matches")
		{
			matchGroup.GET("/:public_id", matchHandler.GetByPublicID)
			matchGroup.POST("/:public_id/forfeit", middleware.AuthMiddleware(cfg.JWTSecret), matchHandler.Forfeit)
			matchGroup.GET("/user/history", middleware.AuthMiddleware(cfg.JWTSecret), matchHandler.History)
		}

		// Leaderboard
		api.GET("/leaderboard", matchHandler.Leaderboard)

		// Submissions & Test Runs
		api.POST("/submissions/run", submissionHandler.RunSample)
		api.POST("/submissions/submit", middleware.AuthMiddleware(cfg.JWTSecret), submissionHandler.SubmitMatch)
	}

	// WebSocket Endpoint
	r.GET("/ws", ws.ServeWS(hub, cfg.JWTSecret))

	serverAddr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("🚀 DevCode API and WebSocket listening at http://localhost%s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}
