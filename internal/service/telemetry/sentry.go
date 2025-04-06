package telemetry

import (
	"os"

	"github.com/JackieLi565/syllabye/internal/config"
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
)

// To initialize Sentry's handler, you need to initialize Sentry itself beforehand

func NewSentryHandler() *sentryhttp.Handler {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv(config.SentryDsn),
	}); err != nil {
		panic("Sentry init failed")
	}

	return sentryhttp.New(sentryhttp.Options{})
}
