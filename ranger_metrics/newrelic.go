package ranger_metrics

import (
	newrelic "github.com/newrelic/go-agent"
	"net/http"
	"fmt"
)

type NewRelic struct {
	Application newrelic.Application
}

func NewNewRelic(appName string, license string) *NewRelic {
	app, err := newrelic.NewApplication(newrelic.NewConfig(
		appName,
		license),
	)

	if err != nil {
		panic(fmt.Errorf("NewRelic error: %s", err))
	}

	return &NewRelic{
		Application: app,
	}
}

func (newRelic *NewRelic) Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		txn := newRelic.Application.StartTransaction(r.URL.Path, w, r)
		defer txn.End()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
