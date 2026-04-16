package goservice

import (
	"fmt"
	"strings"

	golog "github.com/jkaveri/golog/v2"
	"github.com/jkaveri/ramda"

	"github.com/jkaveri/goservice/env"
)

// Config represents the main application configuration
type Config struct {
	// Debug enables debug mode when set to true
	Debug bool

	// DeploymentEnv specifies the deployment environment (e.g. production,
	// staging)
	DeploymentEnv string

	// GRPCServer contains the configuration for the gRPC server
	GRPCServer ServerConfig

	// HTTPServer contains the configuration for the HTTP server
	HTTPServer ServerConfig

	// HealthServer contains the configuration for the health check server
	HealthServer ServerConfig

	// Log contains the logging configuration
	Log LogConfig
}

// ServerConfig represents the configuration for a server instance
type ServerConfig struct {
	// Port defines the network port number for the server to listen on
	Port int
}

// LogConfig holds logging configuration for goconfig loading.
// It is converted to golog.Config in initLogger.
type LogConfig struct {
	// Format is "text" or "json". Defaults to "text".
	Format string
	// Level is "debug", "info", "warn", or "error". Defaults to "debug".
	Level string
	// AddSource adds file:line to each log record when true.
	AddSource bool
}

// toGologConfig converts LogConfig to golog.Config for [golog.InitDefault] /
// [golog.NewLogger].
func (c LogConfig) toGologConfig() golog.Config {
	formatStr := strings.TrimSpace(
		ramda.DefaultFn(c.getDefaultFormat, c.Format),
	)
	level := parseLogLevel(c.Level)

	cfg := golog.Config{
		Level:  level,
		Output: "",
	}

	switch strings.ToLower(formatStr) {
	case "json":
		cfg.Format = golog.FormatJSON
	default:
		cfg.Format = golog.FormatText
	}

	if c.AddSource {
		cfg.EnableSource = true
	}

	return cfg
}

func (c LogConfig) getDefaultFormat() string {
	if env.IsProduction() {
		return "json"
	}

	return "text"
}

// parseLogLevel converts a string level to golog.Level (slog.Level).
func parseLogLevel(s string) golog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return golog.LevelDebug
	case "info":
		return golog.LevelInfo
	case "error":
		return golog.LevelError
	default:
		fmt.Printf(
			"invalid log level: %s, we only support debug, info, error\n",
			s,
		)

		return golog.LevelInfo
	}
}
