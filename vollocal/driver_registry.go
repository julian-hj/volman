package vollocal

import (
	"sync"

	"fmt"

	"github.com/cloudfoundry-incubator/volman/voldriver"
)

type DriverRegistry interface {
	Driver(id string) (voldriver.Driver, bool)
	Activated(id string) (bool, error)
	Activate(id string) error
	Drivers() map[string]voldriver.Driver
	Add(id string, driver voldriver.Driver) error
	Keys() []string
}

type registryEntry struct {
	driver    voldriver.Driver
	activated bool
}

func newRegistryEntry(driver voldriver.Driver) *registryEntry {
	return &registryEntry{driver: driver}
}

type driverRegistry struct {
	sync.RWMutex
	registryEntries map[string]*registryEntry
}

func NewDriverRegistry() DriverRegistry {
	return &driverRegistry{
		registryEntries: make(map[string]*registryEntry),
	}
}

func (d *driverRegistry) Driver(id string) (voldriver.Driver, bool) {
	d.RLock()
	defer d.RUnlock()

	driverEntry, found := d.registryEntries[id]
	if !found {
		return nil, false
	}
	return driverEntry.driver, true
}

func (d *driverRegistry) Drivers() map[string]voldriver.Driver {
	d.RLock()
	defer d.RUnlock()

	driversCopy := map[string]voldriver.Driver{}
	for name, registryEntry := range d.registryEntries {
		driversCopy[name] = registryEntry.driver
	}

	return driversCopy
}

func (d *driverRegistry) Add(id string, driver voldriver.Driver) error {
	d.Lock()
	defer d.Unlock()

	if _, found := d.registryEntries[id]; !found {
		return fmt.Errorf("driver-exists")
	}

	d.registryEntries[id] = newRegistryEntry(driver)
	return nil
}

func (d *driverRegistry) Keys() []string {
	d.Lock()
	defer d.Unlock()

	var keys []string
	for k := range d.registryEntries {
		keys = append(keys, k)
	}

	return keys
}

func (d *driverRegistry) Activated(id string) (bool, error) {
	d.Lock()
	defer d.Unlock()

	driverEntry, found := d.registryEntries[id]
	if !found {
		return false, fmt.Errorf("driver-not-found")
	}

	return driverEntry.activated, nil
}

func (d *driverRegistry) Activate(id string) error {
	d.Lock()
	defer d.Unlock()

	driverEntry, found := d.registryEntries[id]
	if !found {
		return fmt.Errorf("driver-not-found")
	}

	driverEntry.activated = true
	return nil
}
