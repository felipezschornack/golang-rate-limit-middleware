package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/felipezschornack/golang-rate-limit-middleware/internal/limiter/strategy"
)

type RateLimiterMiddleware struct {
	rateLimiter strategy.RateLimiterStrategy
}

func NewRateLimiterMiddleware(rl strategy.RateLimiterStrategy) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{rateLimiter: rl}
}

func (mid *RateLimiterMiddleware) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if mid.rateLimiter.IsBlocked(r) {
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(
				"you have reached the maximum number of requests or actions allowed within a certain time frame",
			)
			return
		}
		next.ServeHTTP(w, r)
	})
}
