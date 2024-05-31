package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/felipezschornack/golang-rate-limit-middleware/internal/limiter"
	"github.com/felipezschornack/golang-rate-limit-middleware/internal/middleware"
	"github.com/go-chi/chi"
)

func main() {
	addr := fmt.Sprintf(":%s", os.Getenv("WEB_SERVER_PORT"))
	http.ListenAndServe(addr, getRouter())
}

func getRouter() *chi.Mux {
	r := chi.NewRouter()

	rateLimiter := limiter.NewRedisRateLimiter()
	middleware := middleware.NewRateLimiterMiddleware(rateLimiter)
	r.Use(middleware.RateLimit)

	r.Get("/", handler)
	return r
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Request not blocked by rate limit!"))
}
