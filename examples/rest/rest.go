package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/stepsisters/kgb"
	"github.com/stepsisters/kgb/pkg/component/kubernetes"
	"github.com/stepsisters/kgb/pkg/component/signal"
	"github.com/stepsisters/kgb/pkg/version"
	"net/http"
	"os"
)

type hello struct {
}

func (h hello) Run() (stop func(), wait func() error, err error) {
	port := os.Getenv("PORT")
	server := &http.Server{Addr: ":" + port, Handler: router()}
	return func() {
			// We got a signal to stop our server
			log.Warnf("%s is shutting down...", version.ApplicationName)
			err := server.Shutdown(context.Background())
			if err != nil {
				log.Errorf("%s is stopped with error, %s", version.ApplicationName, err)
			} else {
				log.Infof("%s is stopped successfully", version.ApplicationName)
			}
		}, func() error {
			return server.ListenAndServe()
		},
		nil
}

func (h hello) IsReady() bool {
	return true
}

// Router register necessary routes and returns an instance of a router.
func router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/welcome", welcome).Methods("GET")
	return r
}

func welcome(w http.ResponseWriter, _ *http.Request) {
	if _, err := fmt.Fprint(w, "Hello! Your request was processed."); err != nil {
		log.Errorf("rest: something went wrong. %s", err)
	}
}

func main() {
	helloComponent := hello{}
	probesPort := os.Getenv("PROBES_PORT")

	components := map[string]kgb.ComponentInterface{
		"k8s-probes":   kubernetes.NewHTTPProbe(helloComponent.IsReady, probesPort),
		"signals-trap": signal.NewTrap(),
		"hello-world":  helloComponent,
	}
	kgb.Run(components)
}