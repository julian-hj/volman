package vollocal

import (
	"os"
	"time"

	"github.com/cloudfoundry-incubator/volman"
	"github.com/cloudfoundry-incubator/volman/voldriver"
	"github.com/onsi/ginkgo"
	"github.com/pivotal-golang/clock"
	"github.com/pivotal-golang/lager"
	"github.com/tedsuo/ifrit"
)

type DriversRegistry struct {
	logger         lager.Logger
	client         volman.Manager
	driversFactory DriverFactory
	scanInterval   time.Duration
	clock          clock.Clock

	DriversMap map[string]voldriver.Driver
}

func NewRegistry(logger lager.Logger, driverFactory DriverFactory, scanInterval time.Duration, clock clock.Clock) *DriversRegistry {
	return &DriversRegistry{
		logger:         logger,
		driversFactory: driverFactory,
		scanInterval:   scanInterval,
		clock:          clock,

		DriversMap: map[string]voldriver.Driver{},
	}
}

func (r *DriversRegistry) RegistryRunner(logger lager.Logger) (ifrit.Runner, error) {
	logger.Info("start")
	defer logger.Info("end")

	return ifrit.RunFunc(func(signals <-chan os.Signal, ready chan<- struct{}) error {
		defer ginkgo.GinkgoRecover()

		interval := r.scanInterval
		timer := r.clock.NewTimer(interval)
		defer timer.Stop()
		for {
			select {
			case <-timer.C():

			case signal := <-signals:
				logger.Info("received-signal", lager.Data{"signal": signal.String()})
			}
			r.SetDrivers(logger)
			close(ready)
			timer.Reset(interval)
		}
		return nil
	}), nil

	return nil, nil
}

func (r *DriversRegistry) SetDrivers(logger lager.Logger) {
	logger = logger.Session("SetDrivers")
	logger.Info("start")
	defer logger.Info("end")

	startime := r.clock.Now()
	logger.Info("set-drivers-startime", lager.Data{"time": startime})

	reg, err := r.driversFactory.Discover(logger)
	if err != nil {
		r.DriversMap = map[string]voldriver.Driver{}
	}

	endtime := r.clock.Now()
	logger.Info("set-drivers-endtime", lager.Data{"time": endtime})

	var driver voldriver.Driver
	for driverName, _ := range reg {
		driver, err = r.driversFactory.Driver(logger, driverName)
		if err != nil {
			r.DriversMap = map[string]voldriver.Driver{}
		}
		r.DriversMap[driverName] = driver
	}
}
