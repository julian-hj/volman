package driverhttp

import (
	"bufio"
	"encoding/json"

	"path/filepath"

	"fmt"

	"github.com/cloudfoundry-incubator/volman/voldriver"
	"github.com/pivotal-golang/lager"

	"github.com/cloudfoundry/gunk/os_wrap"
	"github.com/pivotal-golang/clock"
)

type remoteClient struct {
	httpClient *httpClient
}

func New(url string, clock clock.Clock) voldriver.Driver {
	backoff := NewExponentialBackOff(DefaultBackoffTimeout, clock)

	return &remoteClient{
		httpClient: newHTTPClient(url, voldriver.Routes, backoff),
	}
}

func NewFromPath(driverPath string, osClient os_wrap.Os, clock clock.Clock) (voldriver.Driver, error) {
	extension := filepath.Ext(driverPath)
	switch extension {
	case "sock":
		return New(driverPath, clock), nil
	case "spec":
		return newFromSpec(driverPath, osClient, clock)
	case "json":
		return newFromJSON(driverPath, osClient, clock)
	default:
		return nil, fmt.Errorf("unknown-driver-extension: %s", extension)
	}
}

func newFromSpec(driverPath string, osClient os_wrap.Os, clock clock.Clock) (voldriver.Driver, error) {
	configFile, err := osClient.Open(driverPath)
	if err != nil {
		return nil, err
	}

	addressBytes, _, err := bufio.NewReader(configFile).ReadLine()
	if err != nil {
		return nil, err
	}

	return New(string(addressBytes), clock), nil
}

func newFromJSON(driverPath string, osClient os_wrap.Os, clock clock.Clock) (voldriver.Driver, error) {
	configFile, err := osClient.Open(driverPath)
	if err != nil {
		return nil, err
	}

	var driverJsonSpec voldriver.DriverSpec
	err = json.NewDecoder(configFile).Decode(&driverJsonSpec)
	if err != nil {
		return nil, err
	}

	return New(driverJsonSpec.Address, clock), nil
}

func (r *remoteClient) Activate(logger lager.Logger) voldriver.ActivateResponse {
	logger = logger.Session("activate")
	logger.Info("start")
	defer logger.Info("end")

	var response voldriver.ActivateResponse
	err := r.httpClient.Do(logger, voldriver.ActivateRoute, nil, &response)
	if err != nil {
		logger.Error("failed", err)
		return voldriver.ActivateResponse{Err: err.Error()}
	}

	return response
}

func (r *remoteClient) Create(logger lager.Logger, request voldriver.CreateRequest) voldriver.ErrorResponse {
	logger = logger.Session("create")
	logger.Info("start", lager.Data{"request": request})
	defer logger.Info("end")

	var response voldriver.ErrorResponse
	err := r.httpClient.Do(logger, voldriver.CreateRoute, &request, &response)
	if err != nil {
		logger.Error("failed", err)
		return voldriver.ErrorResponse{Err: err.Error()}
	}

	return response
}

func (r *remoteClient) Mount(logger lager.Logger, request voldriver.MountRequest) voldriver.MountResponse {
	logger = logger.Session("mount")
	logger.Info("start", lager.Data{"request": request})
	defer logger.Info("end")

	var response voldriver.MountResponse
	err := r.httpClient.Do(logger, voldriver.MountRoute, &request, &response)
	if err != nil {
		logger.Error("failed", err)
		return voldriver.MountResponse{Err: err.Error()}
	}

	return response
}

func (r *remoteClient) Unmount(logger lager.Logger, request voldriver.UnmountRequest) voldriver.ErrorResponse {
	logger = logger.Session("unmount")
	logger.Info("start", lager.Data{"request": request})
	defer logger.Info("end")

	var response voldriver.ErrorResponse
	err := r.httpClient.Do(logger, voldriver.UnmountRoute, &request, &response)
	if err != nil {
		logger.Error("failed", err)
		return voldriver.ErrorResponse{Err: err.Error()}
	}

	return response
}

func (r *remoteClient) Remove(logger lager.Logger, request voldriver.RemoveRequest) voldriver.ErrorResponse {
	logger = logger.Session("remove")
	logger.Info("start", lager.Data{"request": request})
	defer logger.Info("end")

	var response voldriver.ErrorResponse
	err := r.httpClient.Do(logger, voldriver.UnmountRoute, &request, &response)
	if err != nil {
		logger.Error("failed-remove", err)
		return voldriver.ErrorResponse{Err: err.Error()}
	}

	return response
}

func (r *remoteClient) Get(logger lager.Logger, request voldriver.GetRequest) voldriver.GetResponse {
	logger = logger.Session("get")
	logger.Info("start", lager.Data{"request": request})
	defer logger.Info("end")

	var response voldriver.GetResponse
	err := r.httpClient.Do(logger, voldriver.UnmountRoute, &request, &response)
	if err != nil {
		logger.Error("failed-remove", err)
		return voldriver.GetResponse{Err: err.Error()}
	}

	return response
}
