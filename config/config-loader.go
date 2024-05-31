package config

import (
	"github.com/spf13/viper"
)

type EnvironmentVariables struct {
	RedisHostname                      string `mapstructure:"REDIS_HOSTNAME"`
	RedisPort                          string `mapstructure:"REDIS_PORT"`
	AccessTokenRateLimitInSeconds      string `mapstructure:"ACCESS_TOKEN_RATE_LIMIT_IN_SECONDS"`
	AccessTokenBlockingWindowInSeconds string `mapstructure:"ACCESS_TOKEN_BLOCKING_WINDOW_IN_SECONDS"`
	IpRateLimitInSeconds               string `mapstructure:"IP_RATE_LIMIT_IN_SECONDS"`
	IpBlockingWindowInSeconds          string `mapstructure:"IP_BLOCKING_WINDOW_IN_SECONDS"`
	WebServerPort                      string `mapstructure:"WEB_SERVER_PORT"`
}

func LoadConfig(path string) (*EnvironmentVariables, error) {
	var cfg *EnvironmentVariables
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
