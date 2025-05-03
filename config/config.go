package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/michaelyusak/go-helper/entity"
	"github.com/sirupsen/logrus"
)

type JwtConfig struct {
	Issuer string `json:"issuer"`
	Key    string `json:"key"`
}

type ServiceConfig struct {
	Port           string          `json:"port"`
	GracefulPeriod entity.Duration `json:"graceful_perion_s"`
	ContextTimeout entity.Duration `json:"context_timeout_s"`
	AllowedOrigins []string        `json:"allowed_origins"`
	MySQL          entity.DBConfig `json:"mysql"`
	Jwt            JwtConfig       `json:"jwt"`
}

func Init(log *logrus.Logger) ServiceConfig {
	configPath := os.Getenv("KREDIT_PLUS_USERS_SERVICE_CONFIG")

	var config ServiceConfig

	configData, err := os.ReadFile(configPath)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": fmt.Sprintf("[config][Init][Read] error: %s", err.Error()),
		}).Fatal("error initiating config file")

		return config
	}

	fmt.Println(string(configData))

	err = json.Unmarshal(configData, &config)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": fmt.Sprintf("[config][Init][Unmarshal] error: %s", err.Error()),
		}).Fatal("error initiating config file")

		return config
	}

	return config
}
