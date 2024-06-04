package limiter

import (
	"strconv"
	"testing"
	"time"

	"github.com/felipezschornack/golang-rate-limit-middleware/config"
	"github.com/felipezschornack/golang-rate-limit-middleware/internal/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type limiterSuite struct {
	suite.Suite
	envVar      *config.EnvironmentVariables
	redisClient *db.Redis
	middleware  *RedisRateLimiter
	ip          string
	token       string
}

func (suite *limiterSuite) SetupSuite() {
	suite.envVar = &config.EnvironmentVariables{
		RedisHostname:                      "localhost",
		RedisPort:                          "6379",
		AccessTokenRateLimitInSeconds:      "10",
		AccessTokenBlockingWindowInSeconds: "3",
		IpRateLimitInSeconds:               "10",
		IpBlockingWindowInSeconds:          "3",
	}
	suite.ip = "127.0.0.1"
	suite.token = uuid.New().String()

	suite.redisClient = db.GetRedisClient(suite.envVar)
	suite.middleware = NewRedisRateLimiter(suite.envVar, suite.redisClient)
}

func (suite *limiterSuite) SetupTest() {
	suite.redisClient.ClearDatabase()
}

func (suite *limiterSuite) TestIpRateLimits() {

	var notBlocked int

	ipRateLimitInSeconds, _ := strconv.Atoi(suite.envVar.IpRateLimitInSeconds)

	totalTimeInSeconds := 11
	var timeout = time.After(time.Second * time.Duration(totalTimeInSeconds))
	blocked := true
	var count int
	var total int
loop:
	for {
		select {
		case <-timeout:
			break loop
		default:
			if !suite.middleware.IsBlocked("", suite.ip) {
				notBlocked++
				if blocked {
					count++
					blocked = false
				}
			} else {
				if !blocked {
					blocked = true
				}
			}
			total++
		}
	}

	assert.LessOrEqual(suite.T(), count*ipRateLimitInSeconds, notBlocked)
}

func (suite *limiterSuite) TestTokenRateLimits() {
	var notBlocked int

	tokenRateLimitInSeconds, _ := strconv.Atoi(suite.envVar.AccessTokenRateLimitInSeconds)

	totalTimeInSeconds := 11
	var timeout = time.After(time.Second * time.Duration(totalTimeInSeconds))
	blocked := true
	var count int
	var total int
loop:
	for {
		select {
		case <-timeout:
			break loop
		default:
			if !suite.middleware.IsBlocked(suite.token, suite.ip) {
				notBlocked++
				if blocked {
					count++
					blocked = false
				}
			} else {
				if !blocked {
					blocked = true
				}
			}
			total++
		}
	}

	assert.LessOrEqual(suite.T(), count*tokenRateLimitInSeconds, notBlocked)
}

func TestWebServerSuite(t *testing.T) {
	suite.Run(t, new(limiterSuite))
}
