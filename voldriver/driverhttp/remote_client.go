package driverhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/cloudfoundry-incubator/cf_http"
	"github.com/cloudfoundry-incubator/volman/voldriver"
	"github.com/pivotal-golang/lager"
	"github.com/tedsuo/rata"
)

type remoteClient struct {
	HttpClient *http.Client
	reqGen     *rata.RequestGenerator
}

func NewRemoteClient(url string) *remoteClient {
	return &remoteClient{
		HttpClient: cf_http.NewClient(),
		reqGen:     rata.NewRequestGenerator(url, voldriver.Routes),
	}
}

func (r *remoteClient) Info(logger lager.Logger) (voldriver.InfoResponse, error) {
	logger = logger.Session("info")
	logger.Info("start")
	defer logger.Info("end")
	return voldriver.InfoResponse{}, nil
}

func (r *remoteClient) Mount(logger lager.Logger, mountRequest voldriver.MountRequest) voldriver.MountResponse {
	volumeID := ""
	logger = logger.Session("remoteclient-mount")
	logger.Info("start")
	defer logger.Info("end")

	sendingJson, err := json.Marshal(mountRequest)
	if err != nil {
		wrappedErr := r.clientError(logger, err, fmt.Sprintf("Error marshalling JSON request %#v", mountRequest))
		return voldriver.MountResponse{Err: wrappedErr.Error()}
	}

	request, err := r.reqGen.CreateRequest(voldriver.MountRoute, nil, bytes.NewReader(sendingJson))

	if err != nil {
		wrappedErr := r.clientError(logger, err, fmt.Sprintf("Error creating request to %s", voldriver.MountRoute))
		return voldriver.MountResponse{Err: wrappedErr.Error()}
	}

	response, err := r.HttpClient.Do(request)
	if err != nil {
		wrappedErr := r.clientError(logger, err, fmt.Sprintf("Error mounting volume %s", volumeID))
		return voldriver.MountResponse{Err: wrappedErr.Error()}
	}

	if response.StatusCode == 500 {
		var remoteError voldriver.Error
		if err := unmarshallJSON(logger, response.Body, &remoteError); err != nil {
			wrappedErr := r.clientError(logger, err, fmt.Sprintf("Error parsing 500 response from %s", voldriver.MountRoute))
			return voldriver.MountResponse{Err: wrappedErr.Error()}
		}
		return voldriver.MountResponse{Err: remoteError.Error()}
	}

	var mountPoint voldriver.MountResponse
	if err := unmarshallJSON(logger, response.Body, &mountPoint); err != nil {
		wrappedErr := r.clientError(logger, err, fmt.Sprintf("Error parsing response from %s", voldriver.MountRoute))

		return voldriver.MountResponse{Err: wrappedErr.Error()}
	}

	return mountPoint
}

func (r *remoteClient) Unmount(logger lager.Logger, unmountRequest voldriver.UnmountRequest) error {
	logger = logger.Session("mount")
	logger.Info("start")
	defer logger.Info("end")

	payload, err := json.Marshal(unmountRequest)
	if err != nil {
		return r.clientError(logger, err, fmt.Sprintf("Error marshalling JSON request %#v", unmountRequest))
	}

	request, err := r.reqGen.CreateRequest(voldriver.UnmountRoute, nil, bytes.NewReader(payload))
	if err != nil {
		return r.clientError(logger, err, fmt.Sprintf("Error creating request to %s", voldriver.UnmountRoute))
	}

	response, err := r.HttpClient.Do(request)
	if err != nil {
		return r.clientError(logger, err, fmt.Sprintf("Error unmounting volume %s", unmountRequest.VolumeId))
	}

	if response.StatusCode == 500 {
		var remoteError voldriver.Error
		if err := unmarshallJSON(logger, response.Body, &remoteError); err != nil {
			return r.clientError(logger, err, fmt.Sprintf("Error parsing 500 response from %s", voldriver.UnmountRoute))
		}
		return remoteError
	}

	return nil
}

func (r *remoteClient) Create(logger lager.Logger, createRequest voldriver.CreateRequest) voldriver.ErrorResponse {
	logger = logger.Session("create")
	logger.Info("start")
	defer logger.Info("end")

	payload, err := json.Marshal(createRequest)
	if err != nil {
		wrappedErr := r.clientError(logger, err, fmt.Sprintf("Error marshalling JSON request %#v", createRequest))
		return voldriver.ErrorResponse{Err: wrappedErr.Error()}
	}

	request, err := r.reqGen.CreateRequest(voldriver.CreateRoute, nil, bytes.NewReader(payload))
	if err != nil {
		wrappedErr := r.clientError(logger, err, fmt.Sprintf("Error creating request to %s", voldriver.CreateRoute))
		return voldriver.ErrorResponse{Err: wrappedErr.Error()}
	}

	response, err := r.HttpClient.Do(request)
	if err != nil {
		wrappedErr := r.clientError(logger, err, fmt.Sprintf("Error creating volume %s", createRequest.Name))
		return voldriver.ErrorResponse{Err: wrappedErr.Error()}
	}

	if response.StatusCode == 500 {
		var remoteError voldriver.ErrorResponse
		if err := unmarshallJSON(logger, response.Body, &remoteError); err != nil {
			wrappedErr := r.clientError(logger, err, fmt.Sprintf("Error parsing 500 response from %s", voldriver.UnmountRoute))
			return voldriver.ErrorResponse{Err: wrappedErr.Error()}
		}
		return remoteError
	}

	return voldriver.ErrorResponse{}
}

func (r *remoteClient) Get(logger lager.Logger, getRequest voldriver.GetRequest) voldriver.GetResponse {
	return voldriver.GetResponse{}
}

func unmarshallJSON(logger lager.Logger, reader io.ReadCloser, jsonResponse interface{}) error {
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		logger.Error("Error in Reading HTTP Response body from remote.", err)
	}
	err = json.Unmarshal(body, jsonResponse)

	return err
}

func (r *remoteClient) clientError(logger lager.Logger, err error, msg string) error {
	logger.Error(msg, err)
	return err
}
