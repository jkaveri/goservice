package goservice

import (
	"context"
	"os"

	golog "github.com/jkaveri/golog/v2"
)

// hostnameEnricher adds the hostname attribute to log records.
type hostnameEnricher struct{}

func newHostnameEnricher() golog.Enricher {
	return &hostnameEnricher{}
}

func (*hostnameEnricher) Enrich(ctx context.Context, b *golog.RecordBuilder) {
	_ = ctx

	name, _ := os.Hostname()
	b.AddAttr(golog.String("hostname", name))
}
