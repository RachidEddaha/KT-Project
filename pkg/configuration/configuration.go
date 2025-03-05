package configuration

import (
	"fmt"
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
	ConfigLogger   `mapstructure:",squash"`
	ConfigEcho     `mapstructure:",squash"`
	ConfigDatabase `mapstructure:",squash"`
	JWTSecret      string `mapstructure:"JWT_SECRET,required=true"`
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

type ConfigDatabase struct {
	Host               string `mapstructure:"DB_HOST,required=true"`
	Port               string `mapstructure:"DB_PORT,required=true"`
	User               string `mapstructure:"DB_USER,required=true"`
	Password           string `mapstructure:"DB_PASSWORD,required=true"`
	Name               string `mapstructure:"DB_NAME,required=true"`
	MaxOpenConnections int    `mapstructure:"DB_MAX_OPEN_CONNECTIONS,default=5"`
	MaxIdleConnections int    `mapstructure:"DB_MAX_IDLE_CONNECTIONS,default=2"`
}

func (c ConfigDatabase) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s search_path=public",
		c.Host, c.Port, c.User, c.Password, c.Name,
	)
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
