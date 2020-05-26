package main

import (
	"context"
	"fmt"
	"github.com/dairlair/kubernetes-ready-application-golang/pkg/component/kubernetes"
	"github.com/dairlair/kubernetes-ready-application-golang/pkg/component/signal"
	"github.com/dairlair/kubernetes-ready-application-golang/pkg/core"
	"github.com/dairlair/kubernetes-ready-application-golang/pkg/version"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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
	fmt.Fprint(w, "Hello! Your request was processed.")
}

func main() {
	helloComponent := hello{}

	components := map[string]core.ComponentInterface{
		"k8s-probes":   kubernetes.NewHTTPProbe(helloComponent.IsReady, "81"),
		"signals-trap": signal.NewTrap(),
		"hello-world":  helloComponent,
	}
	core.Run(components)
}
