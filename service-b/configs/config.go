package configs

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
)

type Conf struct {
	WeatherApiKey string `mapstructure:"WEATHER_API_KEY"`
}

func LoadConfig(path string) (*Conf, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.AddConfigPath(filepath.Join(path, ".."))
	viper.AddConfigPath(filepath.Join(path, "..", ".."))
	viper.AddConfigPath("/")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo de configuração: %w", err)
	}

	var config Conf
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar configurações: %w", err)
	}

	return &config, nil
}
