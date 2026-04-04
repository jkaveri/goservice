package goservice

import (
	golog "github.com/jkaveri/golog/v2"

	"github.com/jkaveri/goservice/check"
)

func initLogger(cfg *Config) {
	check.PanicIfError(golog.InitDefault(cfg.Log.toGologConfig(), newHostnameEnricher()))
}
