package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"smarthome/db"
	"smarthome/handlers"
	"smarthome/messaging"
	"smarthome/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Set up database connection
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/smarthome")
	database, err := db.New(dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer database.Close()

	log.Println("Connected to database successfully")

	// Initialize temperature service
	temperatureAPIURL := getEnv("TEMPERATURE_API_URL", "http://temperature-api:8081")
	temperatureService := services.NewTemperatureService(temperatureAPIURL)
	log.Printf("Temperature service initialized with API URL: %s\n", temperatureAPIURL)

	// Initialize event publisher (интеграция с микросервисами через RabbitMQ)
	var events messaging.Publisher = &messaging.NoopPublisher{}
	rabbitURL := getEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/")
	if publisher, err := connectPublisherWithRetry(rabbitURL, 12, 2*time.Second); err != nil {
		log.Printf("WARNING: RabbitMQ unavailable (%v). Monolith will run without event publishing.", err)
	} else {
		events = publisher
		defer publisher.Close()
	}

	// Initialize router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API routes
	apiRoutes := router.Group("/api/v1")

	// Register sensor routes
	sensorHandler := handlers.NewSensorHandler(database, temperatureService, events)
	sensorHandler.RegisterRoutes(apiRoutes)

	// Start server
	srv := &http.Server{
		Addr:    getEnv("PORT", ":8080"),
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Server starting on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	log.Println("Server exited properly")
}

func connectPublisherWithRetry(url string, attempts int, delay time.Duration) (*messaging.EventPublisher, error) {
	var lastErr error
	for i := 1; i <= attempts; i++ {
		publisher, err := messaging.NewEventPublisher(url)
		if err == nil {
			return publisher, nil
		}
		lastErr = err
		log.Printf("RabbitMQ connect attempt %d/%d failed: %v", i, attempts, err)
		time.Sleep(delay)
	}
	return nil, lastErr
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
