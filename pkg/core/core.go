package core

import (
	"context"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type ComponentInterface interface {
	// It must return two functions, a stop() function to terminate component run at any time, and a
	// wait() function to block until component run fails.
	// If component can not be launched error must be returned
	Run() (stop func(), wait func() error, err error)
}

func Run(components map[string]ComponentInterface) {
	ctx, done := context.WithCancel(context.Background())
	g, gCtx := errgroup.WithContext(ctx)

	for name, component := range components {
		runComponent(g, gCtx, done, name, component)
	}

	err := g.Wait()
	if err != nil {
		log.Errorf("Error group returned error: %s", err)
	}
}

func runComponent(g *errgroup.Group, gCtx context.Context, done context.CancelFunc, name string, component ComponentInterface) {
	g.Go(func() error {
		log.Infof("[%s] launch...", name)
		stop, wait, err := component.Run()
		if err != nil {
			log.Errorf("[%s] component launch failed: %s", name, err)
			return err
		}
		log.Infof("[%s] component launched successfully", name)

		errChan := make(chan error)
		go func() {
			errChan <- wait()
		}()

		select {
		case err = <-errChan:
			log.Errorf("[%s] component error: %s", name, err)
			done()
		case <-gCtx.Done():
			log.Infof("stop [%s] component...", name)
			stop()
			log.Infof("[%s] component stopped successfully", name)
			return gCtx.Err()
		}

		return nil
	})
}