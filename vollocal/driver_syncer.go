package vollocal

import (
	"time"

	"sync"

	"os"

	"github.com/pivotal-golang/clock"
	"github.com/pivotal-golang/lager"
)

type driverSyncer struct {
	sync.RWMutex
	logger         lager.Logger
	scanInterval   time.Duration
	driverFinder   DriverFinder
	driverRegistry DriverRegistry
	clock          clock.Clock
}

func NewDriverSyncer(
	logger lager.Logger,
	scanInterval time.Duration,
	driverFinder DriverFinder,
	driverRegistry DriverRegistry,
	clock clock.Clock,
) *driverSyncer {
	return &driverSyncer{
		logger:         logger,
		scanInterval:   scanInterval,
		driverFinder:   driverFinder,
		driverRegistry: driverRegistry,
		clock:          clock,
	}
}

func (r *driverSyncer) Run(signals <-chan os.Signal, ready chan<- struct{}) error {
	logger := r.logger.Session("sync-drivers")
	logger.Info("start")
	defer logger.Info("end")

	addNewDriversCh := make(chan error, 1)

	timer := r.clock.NewTimer(0)
	defer timer.Stop()

	for {
		select {
		case err := <-addNewDriversCh:
			if err != nil {
				return err
			}
			if ready != nil {
				close(ready)
				ready = nil
			}
			timer.Reset(r.scanInterval)

		case <-timer.C():
			go func() {
				addNewDriversCh <- r.addNewDrivers(logger)
			}()

		case <-signals:
			return nil
		}
	}
}

func (r *driverSyncer) addNewDrivers(logger lager.Logger) error {
	drivers, err := r.driverFinder.Discover(logger)
	if err != nil {
		return err
	}

	for name, driver := range drivers {
		err := r.driverRegistry.Add(name, driver)
		if err != nil {
			// log the error but don't stop
			logger.Error("failed-to-add-driver", err)
		}
	}

	return nil
}
