package volman

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/cloudfoundry-incubator/cf_http"
	"github.com/pivotal-golang/lager"
)

type ListDriversResponse struct {
	Drivers []string `json:"drivers"`
}

type operationType struct {
	Method  string
	Headers map[string]string
}

var (
	usingGet = operationType{"GET", map[string]string{"Accept": "application/json"}}
	usingPos = operationType{"POST", map[string]string{"Accept": "application/json", "Content-Type": "application/json"}}
)

type Client interface {
	ListDrivers(logger lager.Logger) (ListDriversResponse, error)
}

type remoteClient struct {
	HttpClient *http.Client
	URL        string
}

func NewRemoteClient(volmanURL string) *remoteClient {
	return &remoteClient{
		HttpClient: cf_http.NewClient(),
		URL:        volmanURL,
	}
}

func (r *remoteClient) ListDrivers(logger lager.Logger) (ListDriversResponse, error) {
	logger.Session("list-drivers")
	logger.Info("start")

	request := "/v1/drivers"
	response, err := r.Get(logger, request)
	if err != nil {
		logger.Fatal("Error in Listing Drivers", err)
	}
	var drivers ListDriversResponse
	err = AndReturnJsonIn(logger, response, &drivers)

	if err != nil {
		logger.Fatal("Error in Parsing JSON Response of List Drivers", err)
	}
	logger.Info("complete")
	return drivers, err
}

func (r *remoteClient) request(logger lager.Logger, operation operationType, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(operation.Method, r.URL+path, body)
	if err != nil {
		logger.Fatal("Error in creating HTTP Request", err)
	}
	for header, value := range operation.Headers {
		req.Header.Add(header, value)
	}

	response, err := r.HttpClient.Do(req)
	return response, err

}

func (r *remoteClient) Get(logger lager.Logger, fromPath string) (*http.Response, error) {
	return r.request(logger, usingGet, fromPath, nil)
}

func AndReturnJsonIn(logger lager.Logger, response *http.Response, jsonResponse interface{}) error {

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Fatal("Error in Reading HTTP Response body", err)
	}
	err = json.Unmarshal(body, jsonResponse)

	return err
}
