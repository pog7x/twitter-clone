package config

import (
	"fmt"
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Config struct {
	Debug    bool   `mapstructure:"DEBUG"`
	LogLevel string `mapstructure:"LOG_LEVEL"`

	Host string `mapstructure:"HOST"`
	Port int    `mapstructure:"PORT"`

	DatabaseURL string `mapstructure:"DATABASE_URL"`
}

var (
	once    sync.Once
	config  Config
	cfgFile = ".twitter-clone.dev.yml"
)

func init() {
	once.Do(func() {
		initConfig()
	})
}

func initConfig() {
	v := viper.New()
	if cfgFile != "" {
		v.SetConfigFile(cfgFile)

		if err := v.ReadInConfig(); err != nil {
			panic(fmt.Errorf("reading config %s error: %w", v.ConfigFileUsed(), err))
		}
	}

	v.AutomaticEnv()

	envKeysMap := map[string]interface{}{}
	if err := mapstructure.Decode(config, &envKeysMap); err != nil {
		panic(fmt.Errorf("decoding config to mapstructure error: %w", err))
	}

	for k := range envKeysMap {
		if bindErr := v.BindEnv(k); bindErr != nil {
			panic(fmt.Errorf("binding viper env variable '%s' error: %w", k, bindErr))
		}
	}

	if err := v.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("decoding env configuration error: %w", err))
	}
}

func LoadConfig() *Config {
	return &config
}
