package goservice

import (
	"github.com/jkaveri/goconfig"
	"github.com/jkaveri/goservice/check"
)

func loadConfig() Config {
	cfg := Config{
		Debug:         false,
		DeploymentEnv: "dev",
		GRPCServer: ServerConfig{
			Port: 9000,
		},
		HTTPServer: ServerConfig{
			Port: 8080,
		},
		HealthServer: ServerConfig{
			Port: 8081,
		},
		Log: LogConfig{
			Level:     "debug",
			AddSource: true,
		},
	}

	err := goconfig.Load(&cfg)
	check.PanicIfError(err)

	return cfg
}
