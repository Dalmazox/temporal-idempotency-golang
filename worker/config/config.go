package config

import "github.com/spf13/viper"

type Config struct {
	Environment string
	Database    struct {
		ConnectionString string
	}
	Temporal struct {
		Address string
	}
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config

	cfg.Environment = viper.GetString("ENVIRONMENT")
	cfg.Database.ConnectionString = viper.GetString("CONNECTION_STRING")
	cfg.Temporal.Address = viper.GetString("TEMPORAL_ADDRESS")

	return &cfg, nil
}
