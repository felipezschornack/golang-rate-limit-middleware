package strategy

import "net/http"

type RateLimiterStrategy interface {
	IsBlocked(r *http.Request) bool
}
