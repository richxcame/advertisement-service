package main

import (
	"advertisement-service/internal/handler"
	"advertisement-service/internal/repository"
	"advertisement-service/internal/service"
	"advertisement-service/pkg/cache"
	"advertisement-service/pkg/config"
	"advertisement-service/pkg/database"
	"advertisement-service/pkg/middlewares"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	// Load the configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Connect to the database
	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	// Migrate the database
	migrationQueryBytes, err := os.ReadFile("./migrations/init.sql")
	if err != nil {
		log.Fatalf("Could not read migration file: %v", err)
	}
	database.Migrate(db, string(migrationQueryBytes))

	// Create a new advertisement handler
	adRepo := repository.NewAdvertisementRepository(db)
	adService := service.NewService(*adRepo)
	handler.SetService(adService)

	// Create a new Echo instance
	e := echo.New()

	// Create a new Redis cache
	redisClient := cache.NewRedisCache(cfg.Cache)

	// Middlewares
	// Cache middleware
	e.Use(middlewares.CacheMiddleware(redisClient, 5*time.Minute))
	// Invalidate cache middleware
	e.Use(middlewares.InvalidateCacheMiddleware(redisClient))
	// Logging and Recovery middleware
	middlewares.InitLogging(e)
	// Tracing middleware
	cleanupfn := middlewares.InitTracer(e)
	defer cleanupfn()
	// Prometheus middleware
	middlewares.InitMetrics(e)

	// Group routes under "/advertisements"
	adsGroup := e.Group("/advertisements")
	adsGroup.POST("", handler.CreateAdvertisement)
	adsGroup.GET("/:id", handler.GetAdvertisement)
	adsGroup.PUT("/:id", handler.UpdateAdvertisement)
	adsGroup.DELETE("/:id", handler.DeleteAdvertisement)
	adsGroup.GET("", handler.ListAdvertisements)

	// Start server
	go func() {
		if err := e.Start(":" + cfg.Server.Port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}
