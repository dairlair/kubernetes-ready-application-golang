// Signal component provides ability to handle terminating signals from OS.
package signal

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

// Trap implements the ComponentInterface and provides a component which will stop the service when the OS signal
// is received
type Trap struct {
}

func NewTrap() Trap {
	return Trap{}
}

func (ts Trap) Run() (stop func(), wait func() error, err error) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	return func() {
		// Nothing to do
		log.Warnf("trap: close OS signals trap")
	}, func() error {
		sig := <-signalChannel
		log.Warnf("trap: signal received: %s", sig)
		return nil
	},
	nil
}