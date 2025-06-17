package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"weather-api/internal/api"
)

func main() {
	http.HandleFunc("/healthcheck", api.WithCORS(api.Healthcheck))
	http.HandleFunc("/current", api.WithCORS(api.CurrentWeatherHandler))
	http.HandleFunc("/forecast/daily", api.WithCORS(api.ForecastDailyHandler))
	http.HandleFunc("/forecast/hourly", api.WithCORS(api.ForecastHourlyHandler))
	http.HandleFunc("/history", api.WithCORS(api.HistoryHandler))

	addr := os.Getenv("API_URL")
	if addr == "" {
		log.Println("API_URL not set, using default localhost:8080")
		addr = ":8080" 
	}

	fmt.Println("Go API running on ", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
