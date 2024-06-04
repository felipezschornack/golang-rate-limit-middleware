package strategy

type RateLimiterStrategy interface {
	IsBlocked(token string, ipAddress string) bool
}
