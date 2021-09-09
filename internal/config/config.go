package config

import (
	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
)

// Config structure for parse environment variables
type Config struct {
	RedisAddress  string `env:"REDISADDR,required,notEmpty"`
	RedisPassword string `env:"REDISPASS,required,notEmpty"`
	RedisUserName string `env:"REDISUSER,required,notEmpty"`
}

// NewConfig create new config object.
func NewConfig() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.WithFields(log.Fields{
			"handler": "config",
			"action":  "initialize",
		}).Errorf("unable to pars environment variables %v,", err)
	}

	return &cfg
}
