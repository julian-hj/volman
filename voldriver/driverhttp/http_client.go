package driverhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/cf_http"
	"github.com/pivotal-golang/lager"
	"github.com/tedsuo/rata"
)

type httpClient struct {
	backoff BackOff
	client  *http.Client
	reqGen  *rata.RequestGenerator
}

func newHTTPClient(url string, routes rata.Routes, backoff BackOff) *httpClient {
	var client *http.Client

	if strings.HasSuffix(url, ".sock") {
		client = cf_http.NewUnixClient(url)
		url = fmt.Sprintf("unix://%s", url)
	} else {
		client = cf_http.NewClient()
	}

	reqGen := rata.NewRequestGenerator(url, routes)

	return &httpClient{
		backoff: backoff,
		client:  client,
		reqGen:  reqGen,
	}
}

func (r *httpClient) Do(logger lager.Logger, route string, requestBody, responseBody interface{}) error {
	return r.backoff.Retry(logger, func(logger lager.Logger) error {
		logger = logger.Session("do-request")

		requestPayload, err := json.Marshal(requestBody)
		if err != nil {
			logger.Error("failed-marshalling-request", err)
			return err
		}

		request, err := r.reqGen.CreateRequest(route, nil, bytes.NewBuffer(requestPayload))
		if err != nil {
			logger.Error("request-gen-failed", err)
			return err
		}

		response, err := r.client.Do(request)
		if err != nil {
			logger.Error("request-failed", err)
			return err
		}

		err = unmarshallJSON(logger, response.Body, responseBody)
		if err != nil {
			return err
		}

		return nil
	})
}

func unmarshallJSON(logger lager.Logger, reader io.ReadCloser, jsonResponse interface{}) error {
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		logger.Error("failed-reading-json-input", err)
		return err
	}

	err = json.Unmarshal(body, jsonResponse)
	if err != nil {
		logger.Error("failed-unmarshalling-json", err)
		return err
	}

	return nil
}
