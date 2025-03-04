package configuration

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"os"
)

const (
	configName = "config"
	configType = "env"
	pathConfig = "/cmd"
)

type Config struct {
	ConfigLogger `mapstructure:",squash"`
	ConfigEcho   `mapstructure:",squash"`
}

type ConfigLogger struct {
	LogLevel    string `mapstructure:"LOG_LEVEL" default:"info"`
	LogSaveFile bool   `mapstructure:"LOG_SAVE_FILE"`
}

type ConfigEcho struct {
	AllowedOrigins   string `mapstructure:"ALLOWED_ORIGINS" default:"*"`
	AllowCredentials bool   `mapstructure:"ALLOW_CREDENTIALS"`
	AddressEcho      string `mapstructure:"ADDRESS_ECHO" default:":8080"`
}

func LoadConfiguration() (configuration Config, err error) {
	path, err := os.Getwd()
	path = path + pathConfig

	if err != nil {
		return Config{}, err
	}
	viper.SetConfigName(configName)
	viper.AddConfigPath(path)
	viper.SetConfigType(configType)

	viper.AutomaticEnv() // case variables are in the environment

	err = viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	err = viper.Unmarshal(&configuration)
	if err != nil {
		return Config{}, err
	}

	err = validator.New().Struct(&configuration)
	if err != nil {
		return Config{}, err
	}

	return configuration, nil
}
