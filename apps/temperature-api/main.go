package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

// TemperatureResponse matches the contract expected by smart_home TemperatureService.
type TemperatureResponse struct {
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Timestamp   time.Time `json:"timestamp"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	SensorID    string    `json:"sensor_id"`
	SensorType  string    `json:"sensor_type"`
	Description string    `json:"description"`
}

func main() {
	port := getEnv("PORT", "8081")
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Query form used by smart_home: GET /temperature?location=Living%20Room
	mux.HandleFunc("GET /temperature", handleTemperatureByQuery)

	// Path form used by smart_home GetSensors/GetSensorByID: GET /temperature/{sensorID}
	mux.HandleFunc("GET /temperature/{sensorID}", handleTemperatureByID)

	addr := ":" + port
	log.Printf("temperature-api listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func handleTemperatureByQuery(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Query().Get("location")
	sensorID := r.URL.Query().Get("sensorId")
	writeJSON(w, http.StatusOK, buildResponse(location, sensorID))
}

func handleTemperatureByID(w http.ResponseWriter, r *http.Request) {
	sensorID := r.PathValue("sensorID")
	writeJSON(w, http.StatusOK, buildResponse("", sensorID))
}

func buildResponse(location, sensorID string) TemperatureResponse {
	location, sensorID = resolveLocationAndSensorID(location, sensorID)

	// Fresh random reading on every request (requirement: different value each call).
	value := randomTemperature()

	return TemperatureResponse{
		Value:       value,
		Unit:        "°C",
		Timestamp:   time.Now().UTC(),
		Location:    location,
		Status:      "active",
		SensorID:    sensorID,
		SensorType:  "temperature",
		Description: fmt.Sprintf("Simulated temperature for %s", location),
	}
}

// resolveLocationAndSensorID mirrors the mapping from the assignment template.
func resolveLocationAndSensorID(location, sensorID string) (string, string) {
	location = strings.TrimSpace(location)
	sensorID = strings.TrimSpace(sensorID)

	if location == "" {
		switch sensorID {
		case "1":
			location = "Living Room"
		case "2":
			location = "Bedroom"
		case "3":
			location = "Kitchen"
		default:
			location = "Unknown"
		}
	}

	if sensorID == "" {
		switch location {
		case "Living Room":
			sensorID = "1"
		case "Bedroom":
			sensorID = "2"
		case "Kitchen":
			sensorID = "3"
		default:
			sensorID = "0"
		}
	}

	return location, sensorID
}

func randomTemperature() float64 {
	// Comfortable indoor range; one decimal place for readable JSON.
	raw := 15.0 + rand.Float64()*15.0 // [15.0, 30.0)
	return math.Round(raw*10) / 10
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
