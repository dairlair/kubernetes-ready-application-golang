package greeting

import (
	"github.com/stepsisters/kgb/pkg/version"
	log "github.com/sirupsen/logrus"
)

type Greeter struct {
}

func NewGreeter() Greeter {
	return Greeter{}
}

func (g Greeter) Run() (stop func(), wait func() error, err error) {
	return func() {
		}, func() error {
			log.Infof("Starting service %s", version.ApplicationName)
			log.Infof("commit: %s, build time: %s, release: %s", version.Commit, version.BuildTime, version.Release)
			select{ }
		},
		nil
}