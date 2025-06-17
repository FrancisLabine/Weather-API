package api

import (
	"encoding/json"
	"net/http"
	"weather-api/internal/client"
	"strconv"
)

func Healthcheck(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{}{
        "message": "I am healthy!",
        "status":  "success",
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}


// DAILY
func HistoryHandler(w http.ResponseWriter, r *http.Request) {
	lat := r.URL.Query().Get("lat")
    lon := r.URL.Query().Get("lon")
    if lat == "" || lon == "" {
        lat, lon = "45.5", "-73.5" // Montreal
    }

    daysParam := r.URL.Query().Get("days")
    if daysParam == "" {
        daysParam = "1"
    }

    days, err := strconv.Atoi(daysParam)
    if err != nil || days < 1 || days > 30 {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Invalid 'days' parameter. Please use a number between 1 and 30.",
        })
        return
    }

    processedData, err := client.FetchWeatherHistory(lat, lon, days)
    if err != nil {
        http.Error(w, "Failed to fetch weather history", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(processedData)
}

func ForecastDailyHandler(w http.ResponseWriter, r *http.Request) {
    lat := r.URL.Query().Get("lat")
    lon := r.URL.Query().Get("lon")
    if lat == "" || lon == "" {
        lat, lon = "45.5", "-73.5"
    }

    daysParam := r.URL.Query().Get("days")
    days, _ := strconv.Atoi(daysParam)
    if days < 1 || days > 16 {
        days = 7
    }

    data, err := client.FetchDailyForecast(lat, lon, days)
    if err != nil {
        http.Error(w, "Failed to fetch daily forecast", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}

// HOURLY
func ForecastHourlyHandler(w http.ResponseWriter, r *http.Request) {
    lat := r.URL.Query().Get("lat")
    lon := r.URL.Query().Get("lon")
    if lat == "" || lon == "" {
        lat, lon = "45.5", "-73.5" // Montreal
    }

    hoursParam := r.URL.Query().Get("hours")
    hours, _ := strconv.Atoi(hoursParam)
    if hours < 1 || hours > 168 {
        hours = 48
    }

    data, err := client.FetchHourlyForecast(lat, lon, hours)
    if err != nil {
        http.Error(w, "Failed to fetch hourly forecast", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}

func CurrentWeatherHandler(w http.ResponseWriter, r *http.Request) {
    lat := r.URL.Query().Get("lat")
    lon := r.URL.Query().Get("lon")
    if lat == "" || lon == "" {
        lat, lon = "45.5", "-73.5" // Montreal
    }

    data, err := client.FetchCurrentWeather(lat, lon)
    if err != nil {
        http.Error(w, "Failed to fetch current weather", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}
