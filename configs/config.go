package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	WebServerPort string `mapstructure:"WEB_SERVER_PORT"`
	WeatherApiKey string `mapstructure:"WEATHER_API_KEY"`
}

func LoadConfig(path string) (*Config, error) {
	var cfg *Config
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
