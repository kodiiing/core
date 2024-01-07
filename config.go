package main

import (
	"os"

	"dario.cat/mergo"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment" envconfig:"ENVIRONMENT" default:"development"`
	Port        string `yaml:"port" envconfig:"PORT" default:"8080"`
	Databases   struct {
		Host     string `yaml:"host" envconfig:"DB_HOST" default:"localhost"`
		Port     uint16 `yaml:"port" envconfig:"DB_PORT" default:"5432"`
		User     string `yaml:"user" envconfig:"DB_USER" default:"kodiing"`
		Password string `yaml:"password" envconfig:"DB_PASSWORD" default:"VeryStrongPassword"`
		Name     string `yaml:"database" envconfig:"DB_NAME" default:"kodiing"`
	} `yaml:"database"`
	Search struct {
		Host string `yaml:"host" envconfig:"SEARCH_HOST" default:"localhost"`
		Port string `yaml:"port" envconfig:"SEARCH_PORT" default:"8108"`
		Key  string `yaml:"key" envconfig:"SEARCH_KEY" default:""`
	} `yaml:"search"`
	Otel struct {
		ReceiverOtlpGrpcEndpoint string `yaml:"receiver_otlp_grpc_endpoint" envconfig:"OTEL_RECEIVER_OTLP_GRPC_ENDPOINT"`
		ReceiverOtlpHttpEndpoint string `yaml:"receiver_otlp_http_endpoint" envconfig:"OTEL_RECEIVER_OTLP_HTTP_ENDPOINT"`
	} `yaml:"otel"`
}

func GetConfig(configFile string) (Config, error) {
	var configurationFromEnvironment Config
	err := envconfig.Process("", &configurationFromEnvironment)
	if err != nil {
		return configurationFromEnvironment, err
	}

	var configurationFromYaml Config
	if configFile != "" {
		conf, err := os.Open(configFile)
		if err != nil {
			return Config{}, err
		}
		defer func() {
			err := conf.Close()
			if err != nil {
				log.Error().Err(err).Msg("closing configuration file")
			}
		}()

		err = yaml.NewDecoder(conf).Decode(&configurationFromYaml)
		if err != nil {
			return Config{}, err
		}
	}
	//Environment variables set the precedence
	err = mergo.Merge(&configurationFromYaml, configurationFromEnvironment)
	if err != nil {
		return Config{}, err
	}
	return configurationFromYaml, nil
}
