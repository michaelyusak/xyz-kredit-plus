package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/michaelyusak/go-helper/entity"
	"github.com/michaelyusak/go-helper/helper"
	"github.com/sirupsen/logrus"
)

type LocalStorageConfig struct {
	Path string `json:"path"`
}

type ServiceConfig struct {
	Port              string             `json:"port"`
	GracefulPeriod    entity.Duration    `json:"graceful_perion_s"`
	ContextTimeout    entity.Duration    `json:"context_timeout_s"`
	AllowedOrigins    []string           `json:"allowed_origins"`
	MySQL             entity.DBConfig    `json:"mysql"`
	Jwt               helper.JwtConfig   `json:"jwt"`
	Hash              helper.HashConfig  `json:"hash"`
	LocalMediaStorage LocalStorageConfig `json:"local_media_storage"`
	IsEnableSeeding   bool               `json:"is_enable_seeding"`
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

	err = json.Unmarshal(configData, &config)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": fmt.Sprintf("[config][Init][Unmarshal] error: %s", err.Error()),
		}).Fatal("error initiating config file")

		return config
	}

	return config
}
