package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environment     string `mapstructure:"ENVIRONMENT"`
	Host            string `mapstructure:"HOST"`
	Port            string `mapstructure:"PORT"`

	DBUsername      string `mapstructure:"DB_USERNAME"`
	DBPassword      string `mapstructure:"DB_PASSWORD"`
	DBHost          string `mapstructure:"DB_HOST"`
	DBPort          string `mapstructure:"DB_PORT"`
	DBName          string `mapstructure:"DB_DBNAME"`
	DBUrl           string

	DBMigrationPath string `mapstructure:"DB_MIGRATION_PATH"`
	DBRecreate      bool   `mapstructure:"DB_RECREATE"`

	JWTSecret       string        `mapstructure:"JWT_SECRET"`
	JWTDuration     time.Duration `mapstructure:"JWT_DURATION"`
}

func LoadConfig(name string, path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, err
	}
	config.DBUrl = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.DBUsername,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)
	return config, nil
}
