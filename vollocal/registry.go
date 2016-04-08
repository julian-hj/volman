package vollocal

import (
	"os"

	"github.com/cloudfoundry-incubator/cf-lager"
	"github.com/onsi/ginkgo"
	"github.com/tedsuo/ifrit"
)

//go:generate counterfeiter -o ../volmanfakes/fake_driver_factory.go . DriverFactory

type Registry interface {
	SetDrivers(driversPath string) (map[string]string, error)
}

type DriversRegistry struct {
	DriversPath string
}

func NewRegistry(driversPath string) *DriversRegistry {
	return &DriversRegistry{
		DriversPath: driversPath,
	}
}

func NewRegistryRunner(driversPath string) (ifrit.Runner, error) {
	logger, _ := cf_lager.New("DriversRegistry")
	logger.Info("set-drivers-start")
	defer logger.Info("set-drivers-end")
	return ifrit.RunFunc(func(signals <-chan os.Signal, ready chan<- struct{}) error {
		defer ginkgo.GinkgoRecover()

		done := make(chan struct{})
		go func() {
			defer ginkgo.GinkgoRecover()

			SetDrivers(driversPath)
			close(done)
		}()
		close(ready)

		select {
		case <-signals:
		}
		return nil
	}), nil

	return nil, nil
}

func SetDrivers(driversPath string) (map[string]string, error) {
	logger, _ := cf_lager.New("SetDriverRegistry")
	logger.Info("start")
	defer logger.Info("end")
	driverFactory := NewDriverFactory(driversPath)
	reg, err := driverFactory.Discover(logger)
	if err != nil {

		return map[string]string{}, err
	}
	return reg, nil
}
