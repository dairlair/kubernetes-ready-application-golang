package kubernetes

import (
	"context"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// ReadinessChecker is dependency which is used by HTTPProbe struct to check that application is ready to serve.
type ReadinessChecker func () bool

type probe struct {
	readinessChecker ReadinessChecker
}

// HTTPProbe implements the ComponentInterface and provides a component which will handle web request to standard
// kubernetes probes: liveness, readiness and startup probes.
type HTTPProbe struct {
	probe probe
	port string
}

func NewHTTPProbe(readinessChecker ReadinessChecker, port string) HTTPProbe {
	return HTTPProbe{
		probe: probe{
			readinessChecker: readinessChecker,
		},
		port: port,
	}
}

func (p HTTPProbe) Run() (stop func(), wait func() error, err error) {
	server := &http.Server{Addr: ":" + p.port, Handler: p.router()}
	return func() {
			// Handle the stop command
			log.Warn("probes: shutting down...")
			err := server.Shutdown(context.Background())
			if err != nil {
				log.Errorf("probes: stopped with error, %s", err)
			} else {
				log.Infof("probes: stopped successfully")
			}
		}, func() error {
			return server.ListenAndServe()
		},
		nil
}

func (p HTTPProbe) router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/healthz", healthz).Methods("GET")
	r.HandleFunc("/readyz", readyz(p.probe.readinessChecker)).Methods("GET")
	return r
}

// healthz is a liveness probe.
func healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// readyz is a readiness probe.
func readyz(readinessChecker ReadinessChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if readinessChecker() == true {
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		}
	}
}