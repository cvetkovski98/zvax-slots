package config

import (
	"github.com/cvetkovski98/zvax-common/pkg/config"
	"github.com/spf13/viper"
)

type Config struct {
	Redis config.Redis `mapstructure:"db"`
}

func LoadConfig(name string) error {
	viper.AddConfigPath("config")
	viper.SetConfigName(name)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	return viper.ReadInConfig()
}

func GetConfig() *Config {
	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		panic(err)
	}
	return cfg
}
