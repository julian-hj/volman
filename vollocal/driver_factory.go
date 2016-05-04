package vollocal

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/cloudfoundry-incubator/volman/voldriver"
	"github.com/cloudfoundry-incubator/volman/voldriver/driverhttp"
	"github.com/cloudfoundry/gunk/os_wrap"
	"github.com/pivotal-golang/clock"
	"github.com/pivotal-golang/lager"
)

//go:generate counterfeiter -o ../volmanfakes/fake_driver_factory.go . DriverFactory

// DriverFinder is responsible for instantiating remote client implementations of the voldriver.Driver interface.
type DriverFinder interface {
	// Discover will compile a list of drivers from the path list in DriversPath.  If the same driver id is found in
	// multiple directories, it will favor the directory found first in the path.
	// if 2 specs are found within the *same* directory, it will choose .sock files first, then .spec files, then .json
	Discover(logger lager.Logger) (map[string]voldriver.Driver, error)
}

type driverFinder struct {
	driversPath string
	osClient    os_wrap.Os
	clock       clock.Clock
}

func NewDriverFinder(driversPath string, osClient os_wrap.Os, clock clock.Clock) DriverFinder {
	return &driverFinder{
		driversPath: driversPath,
		osClient:    osClient,
		clock:       clock,
	}
}

func (r *driverFinder) Discover(logger lager.Logger) (map[string]voldriver.Driver, error) {
	logger = logger.Session("discover")
	logger.Debug("start")
	logger.Info(fmt.Sprintf("Discovering drivers in %s", r.driversPath))
	defer logger.Debug("end")

	paths := filepath.SplitList(r.driversPath)

	endpoints := make(map[string]voldriver.Driver)
	for _, driverPath := range paths {
		//precedence order: sock -> spec -> json
		spec_types := [3]string{"sock", "spec", "json"}
		for _, spec_type := range spec_types {
			matchingDriverSpecs, err := r.getMatchingDriverSpecs(logger, driverPath, spec_type)

			if err != nil {
				// untestable on linux, does glob work differently on windows???
				return map[string]voldriver.Driver{}, fmt.Errorf("Volman cocd vollovnfigured with an invalid driver path '%s', error occured list files (%s)", driverPath, err.Error())
			}
			if len(matchingDriverSpecs) > 0 {
				logger.Debug("driver-specs", lager.Data{"drivers": matchingDriverSpecs})
				endpoints = r.insertIfNotFound(logger, endpoints, driverPath, matchingDriverSpecs)
			}
		}
	}
	return endpoints, nil
}

func (r *driverFinder) insertIfNotFound(logger lager.Logger, endpoints map[string]voldriver.Driver, driverPath string, specs []string) map[string]voldriver.Driver {
	logger = logger.Session("insert-if-not-found")
	logger.Debug("start")
	defer logger.Debug("end")

	for _, spec := range specs {
		re := regexp.MustCompile("([^/]*/)?([^/]*)\\.(sock|spec|json)$")

		segs2 := re.FindAllStringSubmatch(spec, 1)
		if len(segs2) <= 0 {
			continue
		}
		specName := segs2[0][2]
		logger.Debug("insert-unique-spec", lager.Data{"specname": specName})
		_, ok := endpoints[specName]
		if ok == false {
			driver, err := driverhttp.NewFromPath(driverPath, r.osClient, r.clock)
			if err != nil {
				logger.Error("error-creating-driver", err)
				continue
			}

			endpoints[specName] = driver
		}
	}
	return endpoints
}

func (r *driverFinder) getMatchingDriverSpecs(logger lager.Logger, path string, pattern string) ([]string, error) {
	logger.Debug("binaries", lager.Data{"path": path, "pattern": pattern})
	matchingDriverSpecs, err := filepath.Glob(path + "/*." + pattern)
	if err != nil { // untestable on linux, does glob work differently on windows???
		return nil, fmt.Errorf("Volman configured with an invalid driver path '%s', error occured list files (%s)", path, err.Error())
	}
	return matchingDriverSpecs, nil

}
