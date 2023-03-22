package config

import (
	"backend/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	BindAddr string `env:"BACKEND_BIND_ADDR" env-default:":8080" env-description:"Server bind addr:port"`
}

var appConf *Config

var onceApp sync.Once

func GetAppConfig(l *logging.Logger) *Config {
	onceApp.Do(func() {
		l.Info("Read application config...")
		appConf = &Config{}
		if err := cleanenv.ReadEnv(appConf); err != nil {
			header := "One or more environment variables for application config not found. Supported variables: "
			help, _ := cleanenv.GetDescription(appConf, &header)
			l.Warn(help)
		}
	})
	return appConf
}
