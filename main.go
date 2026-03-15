package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

var (
	requestCount  atomic.Int64
	errorCount    atomic.Int64
	startTime     = time.Now()
	version       = "dev"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/metrics", handleMetrics)
	mux.HandleFunc("/api/users", handleUsers)
	mux.HandleFunc("/api/orders", handleOrders)
	mux.HandleFunc("/api/slow", handleSlow)
	mux.HandleFunc("/api/error", handleError)

	log.Printf("demo-app %s starting on :%s", version, port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	requestCount.Add(1)
	respondJSON(w, 200, map[string]string{
		"service": "demo-app",
		"version": version,
		"status":  "running",
	})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, 200, map[string]interface{}{
		"status": "healthy",
		"uptime": time.Since(startTime).String(),
	})
}

// handleMetrics — Prometheus-совместимый формат
func handleMetrics(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(startTime).Seconds()
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintf(w, "# HELP demo_requests_total Total HTTP requests\n")
	fmt.Fprintf(w, "# TYPE demo_requests_total counter\n")
	fmt.Fprintf(w, "demo_requests_total %d\n", requestCount.Load())
	fmt.Fprintf(w, "# HELP demo_errors_total Total errors\n")
	fmt.Fprintf(w, "# TYPE demo_errors_total counter\n")
	fmt.Fprintf(w, "demo_errors_total %d\n", errorCount.Load())
	fmt.Fprintf(w, "# HELP demo_uptime_seconds Uptime in seconds\n")
	fmt.Fprintf(w, "# TYPE demo_uptime_seconds gauge\n")
	fmt.Fprintf(w, "demo_uptime_seconds %.2f\n", uptime)
}

// handleUsers — имитация получения списка пользователей
func handleUsers(w http.ResponseWriter, r *http.Request) {
	requestCount.Add(1)
	users := []map[string]interface{}{
		{"id": 1, "name": "Alice", "email": "alice@example.com"},
		{"id": 2, "name": "Bob", "email": "bob@example.com"},
		{"id": 3, "name": "Charlie", "email": "charlie@example.com"},
	}
	respondJSON(w, 200, users)
}

// handleOrders — имитация получения заказов с рандомными данными
func handleOrders(w http.ResponseWriter, r *http.Request) {
	requestCount.Add(1)
	orders := make([]map[string]interface{}, 5)
	for i := range orders {
		orders[i] = map[string]interface{}{
			"id":     i + 1,
			"amount": rand.Float64() * 1000,
			"status": []string{"pending", "shipped", "delivered"}[rand.Intn(3)],
		}
	}
	respondJSON(w, 200, orders)
}

// handleSlow — имитация медленного запроса (100-2000ms)
func handleSlow(w http.ResponseWriter, r *http.Request) {
	requestCount.Add(1)
	delay := time.Duration(100+rand.Intn(1900)) * time.Millisecond
	time.Sleep(delay)
	respondJSON(w, 200, map[string]interface{}{
		"message":  "slow response",
		"delay_ms": delay.Milliseconds(),
	})
}

// handleError — имитация ошибки (30% шанс)
func handleError(w http.ResponseWriter, r *http.Request) {
	requestCount.Add(1)
	if rand.Float64() < 0.3 {
		errorCount.Add(1)
		respondJSON(w, 500, map[string]string{"error": "random internal error"})
		return
	}
	respondJSON(w, 200, map[string]string{"message": "ok"})
}

func respondJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
